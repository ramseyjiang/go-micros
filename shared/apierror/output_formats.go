package apierror

import "encoding/json"

// GetExternalSafeAPIError returns a safe APIError object encoded in JSON
func (ae APIError) GetExternalSafeAPIError() *APIError {
	newAE := ae
	newAE.App = ""
	newAE.SourceFile = ""
	newAE.SourceFunc = ""
	newAE.ErrorType = 0
	newAE.OriginErrorMessage = ""
	newAE.Stack = nil
	return &newAE
}

// RequestErrorJSONShort returns a safe APIError object encoded in JSON
func (ae APIError) RequestErrorJSONShort(pretty bool) (int, []byte) {
	newAE := ae.GetExternalSafeAPIError()
	return newAE.RequestErrorJSONFull(pretty)
}

// RequestErrorJSONAuto returns a safe APIError object encoded in JSON
func (ae APIError) RequestErrorJSONAuto() (int, []byte) {
	if DebugMode {
		return ae.RequestErrorJSONFull(!CompactJSON)
	}

	return ae.RequestErrorJSONShort(!CompactJSON)
}

// RequestErrorJSONFull returns a full APIError object encoded in JSON
func (ae APIError) RequestErrorJSONFull(pretty bool) (int, []byte) {
	if pretty {
		if b, err := json.MarshalIndent(ae, "", " "); err != nil || len(b) == 0 {
			return ae.ErrorCode, []byte("{ \"message\": \"Internal API Error\" }")
		} else {
			return ae.ErrorCode, b
		}
	}

	if b, err := json.Marshal(ae); err != nil || len(b) == 0 {
		return ae.ErrorCode, []byte("{ \"message\": \"Internal API Error\" }")
	} else {
		return ae.ErrorCode, b
	}
}

// GetJSONBytes returns the APIError object encoded in JSON
func (ae APIError) GetJSONBytes() ([]byte, error) {
	return json.Marshal(ae)
}

// GetJSONString returns the APIError object encoded in JSON
func (ae APIError) GetJSONString() string {
	jbytes, jerr := json.Marshal(ae)
	if jerr == nil {
		return string(jbytes)
	}
	return ""
}

// GetJSONStringPretty returns the APIError object encoded in JSON (pretty format)
func (ae APIError) GetJSONStringPretty() string {
	jbytes, jerr := json.MarshalIndent(ae, "", " ")
	if jerr == nil {
		return string(jbytes)
	}
	return ""
}
