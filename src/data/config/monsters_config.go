package config

import (
	"encoding/json"
	"os"
)

type MonsterConfig struct {
	MaxHealth int `json:"max_health"`
	Strength  int `json:"strength"`
	Dexterity int `json:"dexterity"`
	Hostility int `json:"hostility"`
}

type MonstersRegistry struct {
	Zombie    MonsterConfig `json:"zombie"`
	Vampire   MonsterConfig `json:"vampire"`
	Ghost     MonsterConfig `json:"ghost"`
	Ogre      MonsterConfig `json:"ogre"`
	SnakeMage MonsterConfig `json:"snake_mage"`
	Mimic     MonsterConfig `json:"mimic"`
}

func LoadMonsterConfigs(filePath string) (*MonstersRegistry, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var registry *MonstersRegistry
	err = json.Unmarshal(fileData, &registry)
	if err != nil {
		return nil, err
	}

	return registry, nil
}
