package app

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"github.com/dan13ram/wpokt-backend/models"
	"gopkg.in/yaml.v2"
)

var (
	Config models.Config
)

func InitConfig(configFile string) {
	var yamlFile, err = ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Error reading config file %q: %s\n", configFile, err.Error())
	}
	err = yaml.Unmarshal(yamlFile, &Config)
	if err != nil {
		log.Fatalf("Error unmarshalling config file %q: %s\n", configFile, err.Error())
	}
}
