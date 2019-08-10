package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// StatusOptions describe possible parameters in the yaml file
type StatusOptions struct {
	Emoji   string
	Message string
}

// Config from yaml file
type Config struct {
	Timezone string
	Crons    map[string]StatusOptions
	Default  StatusOptions
}

// ReadConfig reads the config from a yml file
func ReadConfig(filePath string) (*Config, error) {
	config := Config{}
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(fileContents, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
