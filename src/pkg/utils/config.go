package utils

import (
	"encoding/json"
	"io/ioutil"
)

// DefaultConfigLocation is the default location for config file
const DefaultConfigLocation = "./config/config.json"

type (
	// Config is the model that represents the JSON config for our project
	Config struct {
		DefaultTimezone string
		Port            int
		Mode            string
		Databases       []DatabaseConfig // array since there can be multiple nodes for DB

	}

	// DatabaseConfig defines the config for our DB
	DatabaseConfig struct {
		Host                  string
		Port                  int
		User                  string
		Password              string
		Name                  string
		IdleConnections       int
		OpenConnections       int
		Type                  string
		SamplingRateInSeconds int
	}
)

// LoadConfig reads the JSON file and parses it
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

// GetConfig obtains a single value for the config
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
