package apierror

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io/ioutil"
	"math/rand"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// EmptySlice is used for 404 slice replies
var EmptySlice = make([]NullStruct, 0)

// EmptyObject is used for 404 slice replies
var EmptyObject = NullStruct{}

// NullStruct is used for 404 object replies
type NullStruct struct {
}

var (
	ErrObjectNil      = &APIError{ErrorType: ErrorType_ERROR, ErrorCode: 500, ErrorMessage: "Object is nil"}
	ErrObjectNotFound = &APIError{ErrorType: ErrorType_INFO, ErrorCode: 404, ErrorMessage: "Object not found"}
)

func Is404Error(err error) bool {
	return GetErrorCode(err) == 404
}

func GetErrorCode(err error) int {
	if err == nil {
		return 0
	}

	if err == ErrObjectNotFound {
		return 404
	}

	if ae, ok := err.(*APIError); ok {
		if ae.ErrorCode == 404 || ae.originErr == ErrObjectNotFound {
			return 404
		}
		if ae.ErrorCode <= 0 {
			return 500
		}
		return ae.ErrorCode
	}

	if st, ok := status.FromError(err); ok {
		return GRPCCodeToHTTPCode(st.Code())
	}

	return 500
}

// FromPackedDataBase64 returns the APIError object encoded in Protobuf
func (ae *APIError) FromPackedDataBase64(in string) error {
	if len(in) == 0 {
		return nil
	}
	packedBytes, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		return NewAPIError(err, 500, "", "Unable to decode APIError message")
	}

	if err := ae.FromPackedData(packedBytes); err != nil {
		return NewAPIError(err, 500, "", "")
	}

	return nil
}

// FromPackedDataBase64 returns the APIError object encoded in Protobuf
func (aep *APIErrorProto) FromPackedDataBase64(in string) error {
	if len(in) == 0 {
		return nil
	}
	packedBytes, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		return NewAPIError(err, 500, "", "Unable to decode APIErrorProto message")
	}

	if err := aep.FromPackedData(packedBytes); err != nil {
		return NewAPIError(err, 500, "", "")
	}

	return nil
}

// FromPackedData returns the APIError object encoded in Protobuf
func (ae *APIError) FromPackedData(b []byte) error {
	aep := &APIErrorProto{}
	if err := aep.FromPackedData(b); err != nil {
		return NewAPIError(err, 500, "", "")
	}

	ae.FromAPIErrorProto(aep)

	return nil
}

// FromPackedData returns the APIErrorProto object encoded in Protobuf
func (aep *APIErrorProto) FromPackedData(b []byte) error {
	if aep == nil {
		return nil
	}

	data := bytes.NewReader(b)
	r := flate.NewReader(data)

	uncompressedProtoBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return NewAPIError(err, 500, "", "")
	}

	err = proto.Unmarshal(uncompressedProtoBytes, aep)
	if err != nil {
		return NewAPIError(err, 500, "", "")
	}

	return nil
}

// ToPackedDataBase64 returns the APIError object encoded in Protobuf
func (ae *APIError) ToPackedDataBase64(saveSpace bool) (string, error) {
	packedBytes, err := ae.ToPackedData(saveSpace)
	if err != nil {
		return "", NewAPIError(err, 500, "", "")
	}

	return base64.RawURLEncoding.EncodeToString(packedBytes), nil
}

// ToPackedDataBase64 returns the APIErrorProto object encoded in Protobuf
func (aep *APIErrorProto) ToPackedDataBase64() (string, error) {
	packedBytes, err := aep.ToPackedData()
	if err != nil {
		return "", NewAPIError(err, 500, "", "")
	}

	return base64.RawURLEncoding.EncodeToString(packedBytes), nil
}

// ToPackedData returns the APIErrorProto object encoded in Protobuf
func (ae *APIError) ToPackedData(saveSpace bool) ([]byte, error) {
	aep := APIErrorProto{}
	aep.FromAPIError(ae, saveSpace)
	b, err := aep.ToPackedData()
	if err != nil {
		return nil, NewAPIError(err, 500, "", "")
	}

	return b, nil
}

// ToPackedData returns the APIErrorProto object encoded in Protobuf
func (aep *APIErrorProto) ToPackedData() ([]byte, error) {
	var (
		buf bytes.Buffer
		n   int
		err error
		obj []byte
	)

	protoBytes, err := proto.Marshal(aep)
	if err != nil {
		return nil, NewAPIError(err, 500, "", "")
	}

	buf.Reset()
	fw, err := flate.NewWriter(&buf, 9)
	if err != nil {
		return nil, NewAPIError(err, 500, "", "")
	}
	defer fw.Close()

	// s = time.Now()
	n, err = fw.Write(protoBytes)
	if err != nil {
		return nil, NewAPIError(err, 500, "", "")
	}
	_ = n
	fw.Close() // need to explicitly close and not just flush to prevent EOF error

	obj = buf.Bytes()

	return obj, nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// GetPBBytes returns the APIErrorProto object encoded in Protobuf
func (aep *APIErrorProto) GetPBBytes() ([]byte, error) {
	return proto.Marshal(aep)
}

// GetAPIError returns the APIError object
func (aep *APIErrorProto) GetAPIError() *APIError {
	apiError := APIError{}
	apiError.FromAPIErrorProto(aep)
	return &apiError
}

// FromAPIError loads an APIError into an APIErrorProto
func (aep *APIErrorProto) FromAPIError(apiError *APIError, saveSpace bool) {
	if apiError == nil || aep == nil {
		return
	}
	aep.App = apiError.App
	aep.ErrorCode = int32(apiError.ErrorCode)
	aep.ErrorNumber = int32(apiError.ErrorNumber)
	aep.ErrorField = apiError.ErrorField
	aep.ErrorMessage = apiError.ErrorMessage
	aep.Message = apiError.Message
	aep.SourceFile = apiError.SourceFile
	if !saveSpace {
		aep.ErrorType = apiError.ErrorType
		aep.SourceFunc = apiError.SourceFunc
		aep.OriginErrorMessage = apiError.OriginErrorMessage
	}

	for i := range apiError.traceFrames {
		aep.TraceFrames = append(aep.TraceFrames, &apiError.traceFrames[i])
	}

	if apiError.Stack != nil {
		stackAEP := APIErrorProto{}
		stackAEP.FromAPIError(apiError.Stack, saveSpace)
		aep.Stack = &stackAEP
	}
}

// FromAPIErrorProto loads an APIErrorProto into an APIError
func (ae *APIError) FromAPIErrorProto(aep *APIErrorProto) {
	if ae == nil || aep == nil {
		return
	}
	ae.App = aep.App
	ae.ErrorCode = int(aep.ErrorCode)
	ae.ErrorNumber = int(aep.ErrorNumber)
	ae.ErrorType = aep.ErrorType
	ae.ErrorField = aep.ErrorField
	ae.ErrorMessage = aep.ErrorMessage
	ae.OriginErrorMessage = aep.OriginErrorMessage
	ae.Message = aep.Message
	ae.SourceFile = aep.SourceFile
	ae.SourceFunc = aep.SourceFunc

	for i := range aep.TraceFrames {
		if aep.TraceFrames[i] == nil {
			continue
		}
		ae.traceFrames = append(ae.traceFrames, *aep.TraceFrames[i])
	}

	if aep.Stack != nil {
		stackAE := APIError{}
		stackAE.FromAPIErrorProto(aep.Stack)
		ae.Stack = &stackAE
	}
}
