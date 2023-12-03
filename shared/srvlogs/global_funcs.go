package srvlogs

import (
	"context"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

// Init performs all the base configuration of the GlobalLogger
// This function is typically called when the application is starting up
func Init(applicationName string, syslogHost string, debugMode bool) {
	SetApplicationName(applicationName)
	SetDebugMode(debugMode)
	SetSyslogHost(syslogHost)
}

// SetApplicationName sets the application name on the GlobalLogger
func SetApplicationName(applicationName string) {
	GlobalLogger.SetApplicationName(applicationName)
}

// SetSyslogHost sets the syslog host on the GlobalLogger
func SetSyslogHost(host string) {
	GlobalLogger.SetSyslogHost(host)
}

// SetDebugMode sets the syslog host on the GlobalLogger
func SetDebugMode(debugMode bool) {
	GlobalLogger.SetDebugMode(debugMode)
}

// Add a context to the log entry.
func WithContext(ctx context.Context) *logrus.Entry {
	return GlobalLogger.WithContext(ctx)
}

// Tracef is a global helper / convenience function for accessing the GlobalLogger object
func Tracef(format string, args ...interface{}) {
	GlobalLogger.Tracef(format, args...)
}

// Maximum Number of paths to log
var NumPathsToLog = 1

// How many callers are we going back
var CallersNum = 2

// GetCallerFields returns information about the caller, specifically
// the function name, the filename and line number in the file
func GetCallerFields() logrus.Fields {
	resp := make(map[string]interface{})

	OriginFile := ""
	OriginLine := 0
	OriginFunc := ""

	if pc, file, line, ok := runtime.Caller(CallersNum); ok {
		OriginFile = file
		OriginLine = line
		runtimeFuncPtr := runtime.FuncForPC(pc)
		OriginFunc = runtimeFuncPtr.Name()
	}

	OriginFunc = KeepNumDirs(OriginFunc, NumPathsToLog)
	OriginFile = KeepNumDirs(OriginFile, NumPathsToLog)

	if OriginFile != "" {
		resp["origin_file"] = OriginFile
	}
	if OriginLine > 0 {
		resp["origin_line"] = OriginLine
	}
	if OriginFunc != "" {
		resp["origin_func"] = OriginFunc
	}
	return resp
}

func keepNumDirs(str string, lastn int, startat int) string {
	numFound := strings.Count(str[startat:], "/")
	if numFound > lastn {
		return keepNumDirs(str, lastn, 1+startat+strings.Index(str[startat:], "/"))
	}
	return str[startat:]
}

// KeepNumDirs returns the lastn number of directories in a path+filename combo
func KeepNumDirs(str string, lastn int) string {
	return keepNumDirs(str, lastn, 0)
}

// Debugf is a global helper / convenience function for accessing the GlobalLogger object
func Debugf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Debugf(format, args...)
}

// Infof is a global helper / convenience function for accessing the GlobalLogger object
func Infof(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Infof(format, args...)
}

// Printf is a global helper / convenience function for accessing the GlobalLogger object
func Printf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Printf(format, args...)
}

// Warnf is a global helper / convenience function for accessing the GlobalLogger object
func Warnf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Warnf(format, args...)
}

// Warningf is a global helper / convenience function for accessing the GlobalLogger object
func Warningf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Warningf(format, args...)
}

// Errorf is a global helper / convenience function for accessing the GlobalLogger object
func Errorf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Errorf(format, args...)
}

// Fatalf is a global helper / convenience function for accessing the GlobalLogger object
func Fatalf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Fatalf(format, args...)
}

// Panicf is a global helper / convenience function for accessing the GlobalLogger object
func Panicf(format string, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Panicf(format, args...)
}

// Log is a global helper / convenience function for accessing the GlobalLogger object
func Log(level logrus.Level, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Log(level, args...)
}

// Trace is a global helper / convenience function for accessing the GlobalLogger object
func Trace(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Trace(args...)
}

// Debug is a global helper / convenience function for accessing the GlobalLogger object
func Debug(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Debug(args...)
}

// Info is a global helper / convenience function for accessing the GlobalLogger object
func Info(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Info(args...)
}

// Print is a global helper / convenience function for accessing the GlobalLogger object
func Print(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Print(args...)
}

// Warn is a global helper / convenience function for accessing the GlobalLogger object
func Warn(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Warn(args...)
}

// Warning is a global helper / convenience function for accessing the GlobalLogger object
func Warning(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Warning(args...)
}

// Error is a global helper / convenience function for accessing the GlobalLogger object
func Error(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Error(args...)
}

// Fatal is a global helper / convenience function for accessing the GlobalLogger object
func Fatal(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Fatal(args...)
}

// Panic is a global helper / convenience function for accessing the GlobalLogger object
func Panic(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Panic(args...)
}

// Logln is a global helper / convenience function for accessing the GlobalLogger object
func Logln(level logrus.Level, args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Logln(level, args...)
}

// Traceln is a global helper / convenience function for accessing the GlobalLogger object
func Traceln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Traceln(args...)
}

// Debugln is a global helper / convenience function for accessing the GlobalLogger object
func Debugln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Debugln(args...)
}

// Infoln is a global helper / convenience function for accessing the GlobalLogger object
func Infoln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Infoln(args...)
}

// Println is a global helper / convenience function for accessing the GlobalLogger object
func Println(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Println(args...)
}

// Warnln is a global helper / convenience function for accessing the GlobalLogger object
func Warnln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Warnln(args...)
}

// Warningln is a global helper / convenience function for accessing the GlobalLogger object
func Warningln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Warningln(args...)
}

// Errorln is a global helper / convenience function for accessing the GlobalLogger object
func Errorln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Errorln(args...)
}

// Fatalln is a global helper / convenience function for accessing the GlobalLogger object
func Fatalln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Fatalln(args...)
}

// Panicln is a global helper / convenience function for accessing the GlobalLogger object
func Panicln(args ...interface{}) {
	GlobalLogger.WithFields(GetCallerFields()).Panicln(args...)
}
