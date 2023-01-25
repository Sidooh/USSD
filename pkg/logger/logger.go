package logger

import (
	"USSD.sidooh/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var UssdLog = &log.Logger{
	Out: nil,
}

var ServiceLog = &log.Logger{
	Out: nil,
}

func Init() {
	fmt.Println("Initializing USSD subsystem loggers")

	UssdLog = log.New()
	ServiceLog = log.New()

	env := viper.GetString("APP_ENV")

	if env != "TEST" {
		UssdLog.SetOutput(utils.GetLogFile("ussd-" + time.Now().Format("2006-01-02") + ".log"))
		ServiceLog.SetOutput(utils.GetLogFile("service-" + time.Now().Format("2006-01-02") + ".log"))
	}
}
