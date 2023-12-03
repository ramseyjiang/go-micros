package apierror

import (
	"context"

	"github.com/ramseyjiang/go-micros/shared/srvlog"
)

const (
	// DefaultLogFormatJSONNewLine is used to select the JSON log format
	DefaultLogFormatJSONNewLine = 0
	// DefaultLogFormatLineLog is used to select the line log format (default)
	DefaultLogFormatLineLog = 1
	// DefaultLogFormatJSON is used to select the JSON log formatter
	DefaultLogFormatJSON = 2
)

// APILogger defines an APIError logging interface. This can be used to define your own logging mechanism
type APILogger interface {
	HandleLogMessageWithContext(ctx context.Context, ae *APIError)
}

// APILoggerHandler is a convenience type to avoid having to declare a struct
// to implement the APILogger interface, it can be used like this:
//
//	apierror.Logger = apierror.APILoggerHandler(func(ae *APIError) {
//		// handle the log message
//	})
type APILoggerHandler func(ctx context.Context, ae *APIError)

// // HandleLogMessage implements the Handler interface for the above-mentioned APILoggerHandler convenience type
// func (h APILoggerHandler) HandleLogMessage(ae *APIError) {
// 	if ae != nil {
// 		h(nil, ae)
// 	}
// }

// HandleLogMessageWithContext implements the Handler interface for the above-mentioned APILoggerHandler convenience type
func (h APILoggerHandler) HandleLogMessageWithContext(ctx context.Context, ae *APIError) {
	if ae != nil {
		h(ctx, ae)
	}
}

type DefaultAPILoggerHandler struct{}

func (obj *DefaultAPILoggerHandler) HandleLogMessageWithContext(ctx context.Context, ae *APIError) {
	if obj == nil || ae == nil || srvlog.GlobalLogger == nil {
		return
	}
	fields := ae.GetMapStringInterface()
	switch ae.ErrorType {
	case ErrorType_WARNING:
		srvlog.GlobalLogger.WithFields(fields).Warn(ae.Message)

	case ErrorType_INFO:
		srvlog.GlobalLogger.WithFields(fields).Info(ae.Message)

	case ErrorType_DEBUG:
		srvlog.GlobalLogger.WithFields(fields).Debug(ae.Message)

	default:
		// same as ErrorType_ERROR, ErrorType_UNKNOWN:
		srvlog.GlobalLogger.WithFields(fields).Error(ae.Message)

	}
}
