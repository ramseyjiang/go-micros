package apierror

import (
	"encoding/base64"

	"google.golang.org/protobuf/proto"
)

func (ae *APIError) ToBase64String(saveSpace bool) string {
	aep := APIErrorProto{}
	aep.FromAPIError(ae, saveSpace)
	pb, _ := aep.GetPBBytes()
	return base64.RawURLEncoding.EncodeToString(pb)
}

func FromBase64String(in string) *APIError {
	if len(in) == 0 {
		return nil
	}
	b, err := base64.RawURLEncoding.DecodeString(in)
	if err != nil {
		NewAPIError(err, 500, "", "Unable to decode APIError message")
		return nil
	}

	aep := APIErrorProto{}
	if err := proto.Unmarshal(b, &aep); err != nil {
		return nil
	}
	return aep.GetAPIError()
}
