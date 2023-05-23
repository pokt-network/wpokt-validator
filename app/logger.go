package app

import (
	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	logLevel := Config.Logger.Level

	if logLevel == "debug" {
		log.SetLevel(log.DebugLevel)
	} else if logLevel == "info" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "warn" {
		log.SetLevel(log.WarnLevel)
	}
}
