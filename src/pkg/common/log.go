package common

import (
	"log"
	"os"
)

var debugOn = false

func init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "DEBUG" {
		debugOn = true
	}
}

func LogfDebug(format string, v ...interface{}) {
	if !debugOn {
		return
	}

	LogfInfo(format, v...)
}

func LogfInfo(format string, v ...interface{}) {
	log.Printf(format, v...)
}
