package config

import (
	"encoding/json"
	"os"
)

type ItemConfig struct {
	Name         string `json:"name"`          // Название (для отображения)
	Type         string `json:"type"`          // food, scroll, elixir, weapon, treasure
	StatAffected string `json:"stat_affected"` // health, max_health, strength, dexterity, value
	Amount       int    `json:"amount"`        // Величина эффекта (урон, хил, цена)
	Duration     int    `json:"duration"`      // Длительность в ходах (для эликсиров), 0 = мгновенно/навсегда
}

type ItemsRegistry struct {
	Potions   []ItemConfig `json:"potions"`
	Scrolls   []ItemConfig `json:"scrolls"`
	Food      []ItemConfig `json:"food"`
	Weapons   []ItemConfig `json:"weapons"`
	Treasures []ItemConfig `json:"treasures"`
}

func LoadItemsConfigs(filePath string) (*ItemsRegistry, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var registry ItemsRegistry
	err = json.Unmarshal(fileData, &registry)
	if err != nil {
		return nil, err
	}

	return &registry, nil
}
