package app

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func InitLogger() {
	logLevel := strings.ToLower(Config.Logger.Level)
	log.Debug("[LOGGER] Initializing logger with level: ", logLevel)

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	log.Info("[LOGGER] Logger initialized with level: ", logLevel)
}
