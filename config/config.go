package config

import (
	"encoding/json"
	"lb/common"
	"os"
)

func LoadConfig(file string) (common.Config, error) {
	var config common.Config

	bytes, err := os.ReadFile(file)

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)

	if err != nil {
		return config, err
	}

	return config, nil
}