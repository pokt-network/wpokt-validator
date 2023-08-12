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

func TestInitConfig(t *testing.T) {
	t.Run("Config Initialization Success", func(t *testing.T) {
		// Provide valid config and env files (adjust the paths accordingly)
		configFile := "../config.sample.yml"
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

		// Add assertions as needed
		// For example: assert.NotNil(t, Config.MongoDB)
		//               assert.NotNil(t, Config.Ethereum)
		//               ...
	})

	t.Run("Config Initialization No Config File", func(t *testing.T) {
		// Provide empty config file and valid env file (adjust the paths accordingly)
		configFile := ""
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

		// Add assertions as needed
		// For example: assert.NotNil(t, Config.MongoDB)
		//               assert.NotNil(t, Config.Ethereum)
		//               ...
	})
	// Add more test cases as needed
}

func TestValidateConfig(t *testing.T) {
	t.Run("Valid Configuration", func(t *testing.T) {
		// Create a valid Config struct
		// validConfig := models.Config{
		// 	// Initialize fields as needed
		// }

		configFile := "../config.sample.yml"
		envFile := "../sample.env"

		InitConfig(configFile, envFile)

		// Assign the validConfig to the global Config variable
		// Config = validConfig

		// Call the validateConfig function
		validateConfig()

		// No assertions needed as long as validateConfig doesn't panic or throw an error
	})

	t.Run("Invalid Configuration", func(t *testing.T) {
		// Create an invalid Config struct (missing required fields)
		invalidConfig := models.Config{
			// Initialize fields without required values
		}

		// Assign the invalidConfig to the global Config variable
		Config = invalidConfig

		defer func() { log.StandardLogger().ExitFunc = nil }()
		log.StandardLogger().ExitFunc = func(num int) { panic(fmt.Sprintf("exit %d", num)) }

		// Call the validateConfig function
		assert.Panics(t, func() { validateConfig() }, "validateConfig should panic")

	})
	// Add more test cases as needed
}
