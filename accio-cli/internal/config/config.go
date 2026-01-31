package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	APIEndpoint string `yaml:"apiEndpoint"`
	AuthToken   string `yaml:"authToken"`
	DefaultOrg  string `yaml:"defaultOrg"`
	DefaultEnv  string `yaml:"defaultEnv"`
}

func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".accio", "config.yaml"), nil
}

func LoadConfig() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return &Config{
				APIEndpoint: "http://localhost:8080/api/v1",
				DefaultOrg:  "default-org",
				DefaultEnv:  "dev",
			}, nil
		}
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
