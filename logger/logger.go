package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var UssdLog = &log.Logger{}
var ServiceLog = &log.Logger{}

func Init() {
	fmt.Println("Initializing USSD subsystem loggers")

	// TODO: Ensure logs are rotated daily

	//// Set up default Log
	//filename := "logger/sidooh-" + time.Now().Format("2006-01-02") + ".log"
	//file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.SetOutput(file)

	// Set up USSD Log
	filename := "logger/ussd-" + time.Now().Format("2006-01-02") + ".log"
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
	} else {
		UssdLog.SetOutput(file)
	}

	// Set up Service Log
	filename = "logger/service-" + time.Now().Format("2006-01-02") + ".log"
	file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Error(err)
	} else {
		ServiceLog.SetOutput(file)
	}
}
