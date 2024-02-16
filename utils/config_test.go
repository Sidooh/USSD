package utils

import (
	"github.com/spf13/viper"
	"testing"
)

func TestSetupConfig(t *testing.T) {
	viper.Set("TEST", "test")

	SetupConfig(".")

	test := viper.GetString("TEST")
	if test != "test" {
		t.Errorf("Incorrect env, got: %s, want: %v.", test, "test")
	}

	test = viper.Get("TEST").(string)
	if test != "test" {
		t.Errorf("Incorrect env, got: %s, want: %v.", test, "test")
	}
}
