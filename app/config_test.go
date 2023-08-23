package app

import (
	"fmt"
	"io"
	"testing"

	"github.com/dan13ram/wpokt-validator/models"
	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestReadConfigFromConfigFile(t *testing.T) {
	t.Run("Config File Provided", func(t *testing.T) {
		configFile := "../config.sample.yml"

		read := readConfigFromConfigFile(configFile)

		assert.Equal(t, read, true)
		assert.Equal(t, Config.MongoDB.Database, "mongodb-database")
		assert.Equal(t, Config.MongoDB.TimeoutMillis, int64(2000))
	})

	t.Run("No Config File Provided", func(t *testing.T) {
		configFile := ""

		read := readConfigFromConfigFile(configFile)
		assert.Equal(t, read, false)
	})

	t.Run("Invalid Config File Path", func(t *testing.T) {
		configFile := "../config.sample.invalid.yml"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { readConfigFromConfigFile(configFile) }, "readConfigFromConfigFile should panic")
	})

	t.Run("Invalid Config File Contents", func(t *testing.T) {
		configFile := "../README.md"

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { readConfigFromConfigFile(configFile) }, "readConfigFromConfigFile should panic")
	})

}

func TestInitConfig(t *testing.T) {
	t.Run("Config Initialization Success", func(t *testing.T) {
		configFile := "../config.sample.yml"
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

	})

	t.Run("Config Initialization No Config File", func(t *testing.T) {
		configFile := ""
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

	})
}

func TestValidateConfig(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {

		configFile := "../config.sample.yml"
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

		validateConfig()

	})

	t.Run("Invalid Configuration", func(t *testing.T) {
		invalidConfig := models.Config{}

		Config = invalidConfig

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		assert.Panics(t, func() { validateConfig() }, "validateConfig should panic")

	})
}
