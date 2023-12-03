package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	var (
		cfgString  = New("test-string", "", "test string value")
		cfgInt     = New("test-int", 0, "test int value")
		cfgBool    = New("test-bool", false, "test bool value")
		cfgDefault = New("test-default", "default", "test default value")
	)
	os.Setenv("TEST_STRING", "test")
	os.Setenv("TEST_INT", "1")
	os.Setenv("TEST_BOOL", "true")

	Setup("test", "1.0", "", func() {}, func() {})

	if cfgString.Get() != "test" {
		t.Error("string value not set")
	}
	if cfgInt.Get() != 1 {
		t.Error("int value not set")
	}
	if cfgBool.Get() != true {
		t.Error("bool value not set")
	}
	if cfgDefault.Get() != "default" {
		t.Error("default value not set")
	}
}
