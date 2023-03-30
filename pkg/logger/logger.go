package logger

import (
	"USSD.sidooh/utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var UssdLog = &logrus.Logger{
	Out: nil,
}

var ServiceLog = &logrus.Logger{
	Out: nil,
}

func Init() {
	fmt.Println("Initializing USSD subsystem loggers")

	UssdLog = logrus.New()
	ServiceLog = logrus.New()

	env := viper.GetString("APP_ENV")

	logger := viper.GetString("LOGGER")

	if env != "TEST" {
		if logger == "GCP" {
			UssdLog.SetFormatter(NewGCEFormatter(false))
			ServiceLog.SetFormatter(NewGCEFormatter(false))
		} else {
			UssdLog.SetOutput(utils.GetLogFile("ussd-" + time.Now().Format("2006-01-02") + ".log"))
			ServiceLog.SetOutput(utils.GetLogFile("service-" + time.Now().Format("2006-01-02") + ".log"))
		}
	}
}
