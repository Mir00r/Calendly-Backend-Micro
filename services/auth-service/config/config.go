package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	JWTSecret string `yaml:"secret"`
}

// AppConfig holds the loaded configuration
var AppConfig Config

// LoadConfig reads the config.yaml file and populates AppConfig
func LoadConfig(configPath string) {
	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}
}
