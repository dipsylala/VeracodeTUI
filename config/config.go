package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// VeracodeConfig represents the structure of veracode.yml
type VeracodeConfig struct {
	API struct {
		KeyID     string `yaml:"key-id"`
		KeySecret string `yaml:"key-secret"`
	} `yaml:"api"`
	OAuth struct {
		Enabled bool   `yaml:"enabled"`
		Region  string `yaml:"region"`
	} `yaml:"oauth"`
	Packager map[string]interface{} `yaml:"packager"`
}

// LoadConfig reads and parses the Veracode configuration file
func LoadConfig() (*VeracodeConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".veracode", "veracode.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config VeracodeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate required fields
	if config.API.KeyID == "" || config.API.KeySecret == "" {
		return nil, fmt.Errorf("API key-id and key-secret are required in config file")
	}

	return &config, nil
}

func (c *VeracodeConfig) GetAPICredentials() (keyID, keySecret string) {
	return c.API.KeyID, c.API.KeySecret
}
