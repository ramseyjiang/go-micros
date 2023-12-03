// Package config provides a way to define typed config values, instead of using magic viper strings.
package config

import (
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
