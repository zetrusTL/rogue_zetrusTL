package config

import (
	"encoding/json"
	"os"
)

type PlayerConfig struct {
	MaxHealth int `json:"max_health"`
	Dexterity int `json:"dexterity"`
	Strength  int `json:"strength"`
}

func LoadPlayerConfigs(filePath string) (*PlayerConfig, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var registry *PlayerConfig
	err = json.Unmarshal(fileData, &registry)
	if err != nil {
		return nil, err
	}

	return registry, nil
}
