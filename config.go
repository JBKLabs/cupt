package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Config provides the AWS configuration.
type Config struct {
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
}

// GetConfig reads the JSON AWS configuration file into a Config object.
func GetConfig(path string, config *Config) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("failed to read the AWS configuration file: ", err.Error())
		os.Exit(1)
	}

	json.Unmarshal(raw, &config)
}
