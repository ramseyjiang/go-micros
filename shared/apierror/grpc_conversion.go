package apierror

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ramseyjiang/go-micros/shared/srvlogs/v2"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

// GRPCCodeToHTTPCode converts a GRPC error code to a HTTP code
func GRPCCodeToHTTPCode(grpcCode codes.Code) int {

	switch grpcCode {

	case codes.OK:
		return 0
	// OK is returned on success.

	case codes.Canceled:
		// Canceled indicates the operation was canceled (typically by the caller).
		return http.StatusGone

	case codes.InvalidArgument:
		// InvalidArgument indicates client specified an invalid argument.
		// Note that this differs from FailedPrecondition. It indicates arguments
		// that are problematic regardless of the state of the system
		// (e.g., a malformed file name).
		return http.StatusBadRequest

	case codes.DeadlineExceeded:
		// DeadlineExceeded means operation expired before completion.
		// For operations that change the state of the system, this error may be
		// returned even if the operation has completed successfully. For
		// example, a successful response from a server could have been delayed
		// long enough for the deadline to expire.
		return http.StatusRequestTimeout

	case codes.NotFound:
		// NotFound means some requested entity (e.g., file or directory) was
		// not found.
		return http.StatusNotFound

	case codes.AlreadyExists:
		// AlreadyExists means an attempt to create an entity failed because one
		// already exists.
		return http.StatusConflict

	case codes.PermissionDenied:
		// PermissionDenied indicates the caller does not have permission to
		// execute the specified operation. It must not be used for rejections
		// caused by exhausting some resource (use ResourceExhausted
		// instead for those errors). It must not be
		// used if the caller cannot be identified (use Unauthenticated
		// instead for those errors).
		return http.StatusForbidden

	case codes.ResourceExhausted:
		// ResourceExhausted indicates some resource has been exhausted, perhaps
		// a per-user quota, or perhaps the entire file system is out of space.
		return http.StatusInsufficientStorage
		// return http.StatusTooManyRequests

	case codes.FailedPrecondition:
		// FailedPrecondition indicates operation was rejected because the
		// system is not in a state required for the operation's execution.
		// For example, directory to be deleted may be non-empty, a rmdir
		// operation is applied to a non-directory, etc.
		//
		// A litmus test that may help a service implementor in deciding
		// between FailedPrecondition, Aborted, and Unavailable:
		//  (a) Use Unavailable if the client can retry just the failing call.
		//  (b) Use Aborted if the client should retry at a higher-level
		//      (e.g., restarting a read-modify-write sequence).
		//  (c) Use FailedPrecondition if the client should not retry until
		//      the system state has been explicitly fixed. E.g., if a "rmdir"
		//      fails because the directory is non-empty, FailedPrecondition
		//      should be returned since the client should not retry unless
		//      they have first fixed up the directory by deleting files from it.
		//  (d) Use FailedPrecondition if the client performs conditional
		//      REST Get/Update/Delete on a resource and the resource on the
		//      server does not match the condition. E.g., conflicting
		//      read-modify-write on the same resource.
		return http.StatusPreconditionFailed

	case codes.Aborted:
		// Aborted indicates the operation was aborted, typically due to a
		// concurrency issue like sequencer check failures, transaction aborts,
		// etc.
		//
		// See litmus test above for deciding between FailedPrecondition,
		// Aborted, and Unavailable.
		return http.StatusResetContent

	case codes.OutOfRange:
		// OutOfRange means operation was attempted past the valid range.
		// E.g., seeking or reading past end of file.
		//
		// Unlike InvalidArgument, this error indicates a problem that may
		// be fixed if the system state changes. For example, a 32-bit file
		// system will generate InvalidArgument if asked to read at an
		// offset that is not in the range [0,2^32-1], but it will generate
		// OutOfRange if asked to read from an offset past the current
		// file size.
		//
		// There is a fair bit of overlap between FailedPrecondition and
		// OutOfRange. We recommend using OutOfRange (the more specific
		// error) when it applies so that callers who are iterating through
		// a space can easily look for an OutOfRange error to detect when
		// they are done.
		return http.StatusRequestedRangeNotSatisfiable

	case codes.Unimplemented:
		// Unimplemented indicates operation is not implemented or not
		// supported/enabled in this service.
		return http.StatusNotImplemented

	case codes.Internal:
		// Internal errors. Means some invariants expected by underlying
		// system has been broken. If you see one of these errors,
		// something is very broken.
		return http.StatusInternalServerError

	case codes.Unavailable:
		// Unavailable indicates the service is currently unavailable.
		// This is a most likely a transient condition and may be corrected
		// by retrying with a backoff. Note that it is not always safe to retry
		// non-idempotent operations.
		//
		// See litmus test above for deciding between FailedPrecondition,
		// Aborted, and Unavailable.
		return http.StatusServiceUnavailable

	case codes.DataLoss:
		// DataLoss indicates unrecoverable data loss or corruption.
		return http.StatusTeapot

	case codes.Unauthenticated:
		// Unauthenticated indicates the request does not have valid
		// authentication credentials for the operation.
		return http.StatusUnauthorized

	case codes.Unknown:
		// Unknown error. An example of where this error may be returned is
		// if a Status value received from another address space belongs to
		// an error-space that is not known in this address space. Also,
		// errors raised by APIs that do not return enough error information
		// may be converted to this error.
		return http.StatusUnprocessableEntity

	default:
		return http.StatusSeeOther
	}
}

// HTTPCodeToGRPCCode converts a GRPC error code to a HTTP code
func HTTPCodeToGRPCCode(httpCode int) codes.Code {

	switch httpCode {

	case http.StatusOK:
		// OK is returned on success.
		return codes.OK

	case http.StatusGone:
		// Canceled indicates the operation was canceled (typically by the caller).
		return codes.Canceled

	case http.StatusBadRequest:
		// InvalidArgument indicates client specified an invalid argument.
		// Note that this differs from FailedPrecondition. It indicates arguments
		// that are problematic regardless of the state of the system
		// (e.g., a malformed file name).
		return codes.InvalidArgument

	case http.StatusRequestTimeout:
		// DeadlineExceeded means operation expired before completion.
		// For operations that change the state of the system, this error may be
		// returned even if the operation has completed successfully. For
		// example, a successful response from a server could have been delayed
		// long enough for the deadline to expire.
		return codes.DeadlineExceeded

	case http.StatusNotFound:
		// NotFound means some requested entity (e.g., file or directory) was
		// not found.
		return codes.NotFound

	case http.StatusConflict:
		// AlreadyExists means an attempt to create an entity failed because one
		// already exists.
		return codes.AlreadyExists

	case http.StatusForbidden:
		// PermissionDenied indicates the caller does not have permission to
		// execute the specified operation. It must not be used for rejections
		// caused by exhausting some resource (use ResourceExhausted
		// instead for those errors). It must not be
		// used if the caller cannot be identified (use Unauthenticated
		// instead for those errors).
		return codes.PermissionDenied

	case http.StatusInsufficientStorage:
		// ResourceExhausted indicates some resource has been exhausted, perhaps
		// a per-user quota, or perhaps the entire file system is out of space.
		return codes.ResourceExhausted

	case http.StatusPreconditionFailed:
		// FailedPrecondition indicates operation was rejected because the
		// system is not in a state required for the operation's execution.
		// For example, directory to be deleted may be non-empty, a rmdir
		// operation is applied to a non-directory, etc.
		//
		// A litmus test that may help a service implementor in deciding
		// between FailedPrecondition, Aborted, and Unavailable:
		//  (a) Use Unavailable if the client can retry just the failing call.
		//  (b) Use Aborted if the client should retry at a higher-level
		//      (e.g., restarting a read-modify-write sequence).
		//  (c) Use FailedPrecondition if the client should not retry until
		//      the system state has been explicitly fixed. E.g., if a "rmdir"
		//      fails because the directory is non-empty, FailedPrecondition
		//      should be returned since the client should not retry unless
		//      they have first fixed up the directory by deleting files from it.
		//  (d) Use FailedPrecondition if the client performs conditional
		//      REST Get/Update/Delete on a resource and the resource on the
		//      server does not match the condition. E.g., conflicting
		//      read-modify-write on the same resource.
		return codes.FailedPrecondition

	case http.StatusResetContent:
		// Aborted indicates the operation was aborted, typically due to a
		// concurrency issue like sequencer check failures, transaction aborts,
		// etc.
		//
		// See litmus test above for deciding between FailedPrecondition,
		// Aborted, and Unavailable.
		return codes.Aborted

	case http.StatusRequestedRangeNotSatisfiable:
		// OutOfRange means operation was attempted past the valid range.
		// E.g., seeking or reading past end of file.
		//
		// Unlike InvalidArgument, this error indicates a problem that may
		// be fixed if the system state changes. For example, a 32-bit file
		// system will generate InvalidArgument if asked to read at an
		// offset that is not in the range [0,2^32-1], but it will generate
		// OutOfRange if asked to read from an offset past the current
		// file size.
		//
		// There is a fair bit of overlap between FailedPrecondition and
		// OutOfRange. We recommend using OutOfRange (the more specific
		// error) when it applies so that callers who are iterating through
		// a space can easily look for an OutOfRange error to detect when
		// they are done.
		return codes.OutOfRange

	case http.StatusNotImplemented:
		// Unimplemented indicates operation is not implemented or not
		// supported/enabled in this service.
		return codes.Unimplemented

	case http.StatusInternalServerError:
		// Internal errors. Means some invariants expected by underlying
		// system has been broken. If you see one of these errors,
		// something is very broken.
		return codes.Internal

	case http.StatusServiceUnavailable:
		// Unavailable indicates the service is currently unavailable.
		// This is a most likely a transient condition and may be corrected
		// by retrying with a backoff. Note that it is not always safe to retry
		// non-idempotent operations.
		//
		// See litmus test above for deciding between FailedPrecondition,
		// Aborted, and Unavailable.
		return codes.Unavailable

	case http.StatusTeapot:
		// DataLoss indicates unrecoverable data loss or corruption.
		return codes.DataLoss

	case http.StatusUnauthorized:
		// Unauthenticated indicates the request does not have valid
		// authentication credentials for the operation.
		return codes.Unauthenticated

	// case http.StatusTooManyRequests:
	// 	return codes.ResourceExhausted

	case http.StatusUnprocessableEntity:
		return codes.Unknown

	}

	return codes.Unknown
}

// GetOriginGRPCError returns a new APIError, decoded from GRPC error metadata
func GetOriginGRPCError(err error) *APIError {
	if err == nil {
		return nil
	}
	st := status.Convert(err)
	if st.Code() == codes.OK {
		return nil
	}

	// First look for Base64 proto encoded APIError
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.ErrorInfo:
			if t.Domain == "err" && t.Metadata != nil && t.Metadata["e"] != "" {
				ae := &APIError{}
				decodeErr := ae.FromPackedDataBase64(t.Metadata["e"])
				if decodeErr == nil {
					return ae
				}
				srvlogs.Errorf("GetOriginGRPCError FromPackedDataBase64 ERROR: %v", decodeErr)
			}
			// deprecated - remove
			if t.Domain == "apierror" && t.Metadata != nil && t.Metadata["e"] != "" {
				ae := FromBase64String(t.Metadata["e"])
				if ae != nil {
					srvlogs.Errorf("GetOriginGRPCError FromBase64String ERROR: nil")
					return ae
				}
			}
		}
	}

	// Fallback to legacy JSON encoded APIError - deprecated - remove
	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *errdetails.DebugInfo:
			srvlogs.Debugf("Decoding APIError from DebugInfo")

			ae := APIError{}
			if jsonErr := json.Unmarshal([]byte(t.Detail), &ae); jsonErr == nil {
				return &ae
			} else {
				srvlogs.Errorf("GetOriginGRPCError json.Unmarshal: %v", jsonErr)
			}

			// return nil
		}
	}

	return nil
}

// GRPCError returns a GRPC compatible error
func (ae *APIError) GRPCError() error {
	return ae.GRPCStatus().Err()
}

// GRPCStatus implements the grpc status error interface
// this allows errors returned by a GRPC function to be automatically
// converted into the appropriate GRPC status.
func (ae *APIError) GRPCStatus() *status.Status {
	if ae == nil {
		return nil
	}
	st := status.New(HTTPCodeToGRPCCode(ae.ErrorCode), ae.ErrorMessage)

	var details []protoiface.MessageV1
	if str, encodeErr := ae.ToPackedDataBase64(false); encodeErr == nil {
		if len(str) > 4095 {
			str, encodeErr = ae.ToPackedDataBase64(true)
		}
		ei1 := &errdetails.ErrorInfo{
			Domain: "err",
			Metadata: map[string]string{
				"e": str,
			},
		}
		details = append(details, ei1)
	}

	std, _ := st.WithDetails(details...)

	return std
}

func (ae *APIError) GRPCErrorCtx(ctx context.Context) error {
	GRPCSendContext(ctx)
	return ae.GRPCError()
}

func GRPCSendContext(ctx context.Context) {
	if ctx != nil {
		if md, ok := metadata.FromOutgoingContext(ctx); ok {
			sendErr := grpc.SendHeader(ctx, md)
			if sendErr != nil {
				srvlogs.Warnf("GRPCSendContext grpc.SendHeader(md): %v", sendErr)
			}
		} else {
			sendErr := grpc.SendHeader(ctx, nil)
			if sendErr != nil {
				srvlogs.Warnf("GRPCSendContext grpc.SendHeader(nil): %v", sendErr)
			}
		}
	} else {
		srvlogs.Warnf("GRPCSendContext: no context")
	}
}

// UnaryServerInt should be passed to grpc.NewServer in your main
// function so that all errors are converted to APIErrors even
// if they are not generated by this package or have missed
// being called with GRPCError().
func UnaryServerInt() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {

			var internal *APIError
			switch {
			case errors.As(err, &internal):
				srvlogs.Warn("APIError: ", internal.GetJSONString())
				GRPCSendContext(ctx)
				return nil, internal
			default:
				if grpcErr := GetOriginGRPCError(err); grpcErr != nil {
					srvlogs.Warn("DEFAULT APIError: ", grpcErr.GetJSONString())
					// this is already a GRPC error
					GRPCSendContext(ctx)
					return nil, grpcErr
				}

				wrappedErr := convertWrappedToAPIError(err)
				srvlogs.Warn("untraceable apierror returned by ", info.FullMethod)
				GRPCSendContext(ctx)
				return nil, NewAPIError(wrappedErr, 0, "", "")
			}
		}
		return resp, nil
	}
}

func UnaryInterceptor() grpc.ServerOption {
	return grpc.UnaryInterceptor(UnaryServerInt())
}

// StreamServerInt should be passed to grpc.NewServer in your main
// function so that all errors are converted to APIErrors even
// if they are not generated by this package or have missed
// being called with GRPCError().
func StreamServerInt() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		err := handler(ss.Context(), ss)
		if err != nil {
			var internal *APIError
			switch {
			case errors.As(err, &internal):
				srvlogs.Warn("apierror")
				GRPCSendContext(ss.Context())
				return internal
			default:
				if grpcErr := GetOriginGRPCError(err); grpcErr != nil {
					// this is already a GRPC error
					GRPCSendContext(ss.Context())
					return err
				}

				wrappedErr := convertWrappedToAPIError(err)
				srvlogs.Warn("untraceable apierror returned by ", info.FullMethod)
				GRPCSendContext(ss.Context())
				return NewAPIErrorWithContext(ss.Context(), wrappedErr, 500, "", "").GRPCError()
			}
		}
		return nil
	}
}

func StreamInterceptor() grpc.ServerOption {
	return grpc.StreamInterceptor(StreamServerInt())
}
