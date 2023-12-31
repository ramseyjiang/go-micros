package srvlog

import "github.com/ramseyjiang/go-micros/shared/srvlog/logruswrapper"

// GlobalLogger is used by the Global Logging functions
var GlobalLogger = NewGlobalSRVLogger()

// NewGlobalSRVLogger creates and return a new Logger
func NewGlobalSRVLogger() *logruswrapper.LogrusWrapper {
	return logruswrapper.NewLogrusWrapper()
}
