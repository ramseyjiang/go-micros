// Package config provides a way to define typed config values, instead of using magic viper strings.
package config

import (
	"os"

	"github.com/ramseyjiang/go-micros/shared/apierror"
	"github.com/ramseyjiang/go-micros/shared/srvlogs/v2"
	"github.com/ramseyjiang/go-micros/shared/viperconf/v2"
	"github.com/spf13/viper"
)

// Value is a configuration option.
type Value[T any] struct {
	name string
}

// Get returns the configured value.
func (v Value[T]) Get() T {
	var zero T
	viper.UnmarshalKey(v.name, &zero)
	return zero
}

// New creates a new configuration option.
func New[T any](name string, defaultsTo T, description string) Value[T] {
	viperconf.New(name, defaultsTo, description)
	return Value[T]{name}
}

func Setup(applicationName string, version string, buildStamp string, initViperPFlags viperconf.ViperSetParamsFunc, runWhenConfigChanges viperconf.ViperSetParamsFunc) {
	if configErr := viperconf.SetupViperV2(nil, applicationName, initViperPFlags, runWhenConfigChanges); configErr != nil {
		srvlogs.Errorf("Configuration ERROR: %v", configErr)
		os.Exit(1)
	}

	srvlogs.Init(applicationName, viper.GetString("syslog-host"), viper.GetBool("debug"))
	apierror.ApplicationName = applicationName + "/" + version
	apierror.DebugMode = viper.GetBool("debug")

	srvlogs.Infof("%s - Version: %s, Build: %s", applicationName, version, buildStamp)

	if viper.GetBool("debug") {
		viperconf.ShowConfig()
	}

	// strhelpers.TimeZoneInit()
}
