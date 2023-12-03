package apierror

func (ae *APIError) GetMapStringInterface() map[string]interface{} {
	fields := make(map[string]interface{})

	if ae.ErrorCode != 0 {
		fields["err_code"] = ae.ErrorCode
	}
	if ae.ErrorNumber != 0 {
		fields["err_number"] = ae.ErrorNumber
	}
	if ae.ErrorCode != 0 {
		fields["err_field"] = ae.ErrorField
	}

	fields["err_message"] = ae.ErrorMessage

	if ae.OriginErrorMessage != "" {
		fields["origin_err_message"] = ae.OriginErrorMessage
	}

	if ae.App != "" {
		fields["app"] = ae.App
	}
	if ae.SourceFile != "" {
		fields["source_file"] = ae.SourceFile
	}
	if ae.SourceFunc != "" {
		fields["source_func"] = ae.SourceFunc
	}

	if ae.ErrorType > 0 {
		fields["err_type"] = ErrorType_name[int32(ae.ErrorType)]
	}

	if ae.Message != "" {
		fields["message"] = ae.Message
	}

	if ae.Stack != nil {
		fields["stack"] = ae.Stack.GetMapStringInterface()
	}

	return fields
}

// Error implements the standard error handler interface
func (ae *APIError) Error() string {
	if ae == nil || ae.isNil {
		return "<nil>"
	}
	return ae.ErrorMessage
}

// GetErrorCode gets the ErrorCode
func (ae *APIError) GetErrorCode() int {
	return ae.ErrorCode
}

// GetErrorMessage gets the ErrorMessage
func (ae *APIError) GetErrorMessage() string {
	return ae.ErrorMessage
}

// GetErrorField gets the ErrorField
func (ae *APIError) GetErrorField() string {
	return ae.ErrorField
}

// GetErrorNumber gets the ErrorNumber
func (ae *APIError) GetErrorNumber() int {
	return ae.ErrorNumber
}

// GetTraceFrames gets the TraceFrames
func (ae *APIError) GetTraceFrames() []TraceFrame {
	return ae.traceFrames
}
