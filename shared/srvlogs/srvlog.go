package srvlogs

import "github.com/ramseyjiang/micros/shared/srvlogs/v2/logruswrapper"

// GlobalLogger is used by the Global Logging functions
var GlobalLogger = NewGlobalSRVLogger()

// NewGlobalSRVLogger creates and return a new Logger
func NewGlobalSRVLogger() *logruswrapper.LogrusWrapper {
	return logruswrapper.NewLogrusWrapper()
}
