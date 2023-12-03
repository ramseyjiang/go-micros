package logruswrapper

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type logContextData struct {
	// ServiceName describes the application/microservice name
	ServiceName string
}

type contextDataHook struct {
	data *logContextData
}

func (ecdh contextDataHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel,
		logrus.PanicLevel,
	}
}

func (ecdh contextDataHook) Fire(e *logrus.Entry) error {
	if appVal, appExists := e.Data["app"]; appExists && appVal != "" && ecdh.data != nil && ecdh.data.ServiceName != "" {
		// app name exists, leave it
	} else if ecdh.data != nil && ecdh.data.ServiceName != "" {
		e.Data["app"] = ecdh.data.ServiceName
	}
	if e.Logger.GetLevel() != logrus.DebugLevel {
		delete(e.Data, "origin_file")
		delete(e.Data, "origin_line")
		// delete(e.Data, "origin_func")
	}
	if e.Message == "" {
		if errmsg, exists := e.Data["message"]; exists {
			e.Message = fmt.Sprintf("%v", errmsg)
		}
	}
	if e.Message == "" {
		if errmsg, exists := e.Data["err_message"]; exists {
			e.Message = fmt.Sprintf("%v", errmsg)
		}
	}
	return nil
}
