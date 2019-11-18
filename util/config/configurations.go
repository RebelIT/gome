package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var App *Configuration

func LoadConfiguration(configurationFile string) error {
	log.Printf("INFO: LoadConfiguration, %s\n", configurationFile)
	App := &Configuration{}

	data, err := ioutil.ReadFile(configurationFile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &App); err != nil {
		return err
	}

	log.Printf("INFO: LoadConfiguration, done\n")
	return nil
}
