package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host   string `json:"host"`
	Port   string `json:"port"`
	DbFile string `json:"dbFile"`
}

func NewConfig() (*Config, error) {
	body, err := os.ReadFile("./config/config.json")
	if err != nil {
		return nil, err
	}
	newConfig := &Config{}
	if err := json.Unmarshal(body, newConfig); err != nil {
		return nil, err
	}
	return newConfig, nil
}
