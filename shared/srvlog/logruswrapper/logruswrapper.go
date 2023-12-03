package logruswrapper

import (
	"os"
	"runtime"
	"strings"

	srslog "github.com/RackSec/srslog"
	"github.com/sirupsen/logrus"
)

// LogrusWrapper wraps a logrus wrapper with our own context and Syslog information
type LogrusWrapper struct {
	*logrus.Logger
	contextData logContextData
	SyslogHost  string
}

// NewLogrusWrapper initializes the standard logger
func NewLogrusWrapper() *LogrusWrapper {
	newLogrusLogger := logrus.New()
	newLogger := LogrusWrapper{
		Logger: newLogrusLogger,
	}
	newLogger.SetStructuredLogging(false)
	// newLogger.ReportCaller = true
	newLogger.AddHook(contextDataHook{data: &newLogger.contextData})

	// Enable info log level by default
	newLogger.SetLevel(logrus.InfoLevel)

	return &newLogger
}

func (l *LogrusWrapper) SetApplicationName(name string) {
	l.contextData.ServiceName = name
}

func (l *LogrusWrapper) SetDebugMode(debugMode bool) {
	if debugMode {
		l.SetLevel(logrus.DebugLevel)
	} else {
		l.SetLevel(logrus.InfoLevel)
	}
}

func (l *LogrusWrapper) SetStructuredLogging(structured bool) {
	if structured {
		// l.Formatter = &logrus.JSONFormatter{
		// 	TimestampFormat: time.RFC3339Nano,
		// }
		l.Formatter = &IQJSONFormatter{}
		return
	}
	l.Formatter = &IQTextFormatter{
		UseColour: checkIfTerminal(l.Logger.Out) && (runtime.GOOS != "windows"),
	}
}

func (l *LogrusWrapper) SetSyslogHost(newhost string) {
	if newhost != "" && strings.Index(newhost, ":") == -1 {
		// make sure we have a (UDP) port in the host definition
		newhost = newhost + ":514"
	}

	if l.SyslogHost == newhost {
		// no change
		l.Debugf("Syslog host not changed to (%s) - already set to that", newhost)
		return
	}

	if newhost == "" {
		l.Infof("Log output changed to StdErr", newhost)
		l.SetOutput(os.Stderr)

		// Disable structured logging by default when sending to console/StdErr
		l.SetStructuredLogging(false)

		return
	}

	newSyslog, syslogErr := srslog.Dial("udp", newhost, srslog.LOG_DAEMON|srslog.LOG_INFO, l.contextData.ServiceName)
	if syslogErr == nil && newSyslog != nil {
		newSyslog.SetFormatter(srslog.RFC3164Formatter)
		l.SetOutput(newSyslog)

		// Enable structured logging by default when sending to Syslog
		l.SetStructuredLogging(true)

		if l.contextData.ServiceName != "" {
			l.Debugf("Log output for (%s) set to syslog to (%s)", l.contextData.ServiceName, newhost)
		} else {
			l.Debugf("Log output set to syslog to (%s)", newhost)
		}
	} else {
		l.Errorf("ERROR: Unable to init syslog to (%s): %v", newhost, syslogErr)
		os.Exit(1)
	}
}
