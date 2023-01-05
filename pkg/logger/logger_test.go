package logger

import (
	"testing"
)

func TestInit(t *testing.T) {

	if UssdLog.Out != nil || ServiceLog.Out != nil {
		t.Errorf("Init() = %v, %v; want nil", UssdLog.Out, ServiceLog.Out)
	}

	Init()

	if UssdLog.Out == nil || ServiceLog.Out == nil {
		t.Errorf("Init() = %v, %v; want values", UssdLog.Out, ServiceLog.Out)
	}

	if UssdLog.Level.String() != "info" && UssdLog.Level.String() != "info" {
		t.Errorf("Init() = %v, %v; want info", UssdLog.Level, ServiceLog.Level)
	}

}
