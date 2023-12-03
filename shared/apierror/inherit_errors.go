package apierror

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func (ae *APIError) ApplyFieldInheritance() {
	// START decoding field error from github.com/mwitkow/go-proto-validators/validator.proto
	if ae.ErrorField == "" {
		if strings.HasPrefix(ae.ErrorMessage, "invalid field ") {
			prefix := ae.ErrorMessage[14:]
			colonOffset := strings.Index(prefix, ":")
			if colonOffset > 0 {
				if len(prefix) > colonOffset+2 {
					ae.ErrorField = prefix[:colonOffset]
					ae.ErrorMessage = prefix[colonOffset+2:]
					if ae.ErrorCode == 0 {
						ae.ErrorCode = 400 // Bad Input
					}
				}
			}
		}
	}
	if ae.ErrorField == "" {
		oeMsg := ae.OriginErrorMessage
		if ae.OriginErrorMessage == "" && ae.originErr != nil {
			oeMsg = ae.originErr.Error()
		}
		if strings.HasPrefix(oeMsg, "invalid field ") {
			prefix := oeMsg[14:]
			colonOffset := strings.Index(prefix, ":")
			if colonOffset > 0 {
				if len(prefix) > colonOffset+2 {
					ae.ErrorField = prefix[:colonOffset]
					ae.ErrorMessage = prefix[colonOffset+2:]
					if ae.ErrorCode == 0 {
						ae.ErrorCode = 400 // Bad Input
					}
				}
			}
		}
	}
	// END decoding field error from github.com/mwitkow/go-proto-validators/validator.proto

	// Inherit Error Messages from the OriginErrorMessage, then OriginErr, and finally the stack
	if ae.ErrorMessage == "" && ae.OriginErrorMessage != "" {
		ae.ErrorMessage = ae.OriginErrorMessage
	}
	if ae.ErrorMessage == "" && ae.originErr != nil {
		ae.ErrorMessage = ae.originErr.Error()
	}
	if ae.ErrorMessage == "" && ae.Stack != nil {
		ae.ErrorMessage = ae.Stack.ErrorMessage
		// If we are inheriting the ErrorMessage, we can (optionally) inherit the ErrorField as well
		if ae.ErrorField == "" {
			ae.ErrorField = ae.Stack.ErrorField
		}
	}

	// Inherit Error Codes from the Stack (regardless)
	if ae.ErrorCode == 0 && ae.Stack != nil {
		ae.ErrorCode = ae.Stack.ErrorCode
	}
	if ae.ErrorCode == 0 && ae.ErrorType == ErrorType_ERROR {
		ae.ErrorCode = 500
	}

	// Inherit Error Numbers from the Stack (regardless)
	if ae.ErrorNumber == 0 && ae.Stack != nil {
		ae.ErrorNumber = ae.Stack.ErrorNumber
	}
}

// ApplyForeignError applies a foreign error to this APIError object (must return ae if successful, or nil if its not actually an error)
func (ae *APIError) ApplyForeignError(err error, errcode int, errfield string) error {
	if err == nil {
		return ae
	}

	if e, isAPIErrorType := err.(*APIError); isAPIErrorType && e != nil {
		ae.Stack = e
		ae.originErr = e.originErr
		ae.OriginErrorMessage = e.OriginErrorMessage
		ae.ImportChildFields(e)
		return ae
	} else if e, isAPIErrorJSONType := err.(apiErrorInspect); isAPIErrorJSONType && e != nil {
		// This code is more "inter-version compatible"
		if jsonBytes, _ := e.GetJSONBytes(); len(jsonBytes) > 0 {
			stackObj := APIError{}
			if jsonErr := json.Unmarshal(jsonBytes, &stackObj); jsonErr == nil {
				ae.Stack = &stackObj
				ae.originErr = nil
				ae.OriginErrorMessage = ""
				ae.ImportChildFields(&stackObj)
				// srvlog.Debugf("ApplyForeignError: apiErrorInspect")
				return ae
			} else {
				log.Printf("apierror.ApplyForeignError JSON ERROR: %v", jsonErr)
			}
		}
	} else if e, isAPIErrorPBType := err.(apiErrorProtoInspect); isAPIErrorPBType && e != nil {
		// This code is more "inter-version compatible"
		if pbBytes, _ := e.GetPBBytes(); len(pbBytes) > 0 {
			stackObj := APIErrorProto{}
			if jsonErr := proto.Unmarshal(pbBytes, &stackObj); jsonErr == nil {
				ae.Stack = stackObj.GetAPIError()
				ae.originErr = nil
				ae.OriginErrorMessage = ""

				// Convert APIErrorProto to APIError
				tmpAPIError := APIError{}
				tmpAPIError.FromAPIErrorProto(&stackObj)
				ae.ImportChildFields(&tmpAPIError)

				return ae
			} else {
				log.Printf("apierror.ApplyForeignError PB ERROR: %v", jsonErr)
			}
		}
	} else if e, isAPIErrorObjectType := err.(apiErrorObjectInspect); isAPIErrorObjectType && e != nil {
		// This code is more "inter-version compatible"
		if newAE := e.GetAPIError(); newAE != nil {
			ae.Stack = newAE
			ae.originErr = nil
			ae.OriginErrorMessage = ""
			ae.ImportChildFields(newAE)
			return ae
		}
	} else {
		if e, isAPIErrorType := err.(*APIError); isAPIErrorType && e != nil {
			ae.Stack = e
			ae.originErr = nil
			ae.OriginErrorMessage = ""
			ae.ImportChildFields(e)
			return ae
		}
	}

	if echo4Err, ok := err.(echoV4HTTPError); ok && echo4Err != nil {
		if echo4Err.Unwrap() != nil {
			err = echo4Err.Unwrap() // promote the underlying internal error (it might be a json or validation error)
			ae.originErr = echo4Err.Unwrap()
			if ae.originErr != nil {
				ae.OriginErrorMessage = ae.originErr.Error()
			}
		}
	}

	// Check for JSON validation error
	if jsonErr, ok := (err).(*json.UnmarshalTypeError); ok {
		ae.originErr = jsonErr
		// ae.OriginErrorMessage = jsonErr.Error()
		ae.ErrorField = jsonErr.Field
		ae.ErrorMessage = fmt.Sprintf("Invalid type (%s) for (%s), should be (%s)", jsonErr.Value, jsonErr.Field, jsonErr.Type.String())
		return ae
	}

	// GRPC errors are in the format: "rpc error: code = %s"
	if statusErr, isStatusError := status.FromError(err); isStatusError && statusErr != nil {
		if statusErr.Code() == codes.OK {
			return nil // this is not actually an error
		}

		originErr := GetOriginGRPCError(err)
		if originErr != nil {
			// This is an APIerror imported from GRPC (use standard inheritance)
			ae.Stack = originErr
			ae.originErr = nil
			ae.OriginErrorMessage = ""

			return ae
		}

		// this is a GRPC error, but not an embedded APIError
		if ae.ErrorCode <= 0 && ae.Stack != nil && ae.Stack.ErrorCode > 0 {
			ae.ErrorCode = ae.Stack.ErrorCode
		}
		if ae.ErrorCode <= 0 {
			ae.ErrorCode = GRPCCodeToHTTPCode(statusErr.Code())
		}
		if ae.Stack == nil && ae.OriginErrorMessage == "" {
			ae.OriginErrorMessage = statusErr.Message()
		}

		return ae
	}

	rv := reflect.ValueOf(err)
	if rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		p := rv.Pointer()
		if p != 0 {
			ae.OriginErrorMessage = err.Error()
		} else {
			// nil pointer
			if (errcode == 0 || errcode == -1) && errfield == "" && ae.ErrorMessage == "" {
				return nil
				// ae.isNil = true
				// return &ae
			}
		}
	}

	return ae
}

func (ae *APIError) ImportChildFields(child *APIError) {
	if child == nil {
		return
	}
	if ae.ErrorMessage == "" {
		ae.ErrorMessage = child.ErrorMessage
		if ae.ErrorField == "" {
			ae.ErrorField = child.ErrorField
		}
	}
	if ae.ErrorCode == 0 {
		ae.ErrorCode = child.ErrorCode
	}
}

func GetAnyErrorMessage(err error) string {
	if err == nil {
		return "<nil>"
	}

	ae := NewAPIDebug(err, 0, "", "")
	if ae == nil || ae.ErrorMessage == "" {
		return err.Error()
	}

	return ae.ErrorMessage
}
