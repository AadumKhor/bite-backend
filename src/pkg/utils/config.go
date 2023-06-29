package utils

import (
	"encoding/json"
	"io/ioutil"
)

const DefaultConfigLocation = "./config/config.json"

type (
	Config struct {
		DefaultTimezone string
		Port            int
		Mode            string
		Databases       []DatabaseConfig // array since there can be multiple nodes for DB

	}

	DatabaseConfig struct {
		Host                  string
		Port                  string
		User                  string
		Password              string
		Name                  string
		IdleConnections       int
		OpenConnections       int
		Type                  string
		SamplingRateInSeconds int
	}
)

func (c *Config) LoadConfig() error {
	defaultConfig, err := ioutil.ReadFile(DefaultConfigLocation)
	if err != nil {
		return err
	}

	err = json.Unmarshal(defaultConfig, c)
	if err != nil {
		return err
	}

	return nil
}

var config *Config

func GetConfig() (*Config, error) {
	// singleton config check
	if config == nil {
		config = &Config{}
		err := config.LoadConfig()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}
