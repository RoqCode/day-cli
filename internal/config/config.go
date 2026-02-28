package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Ai struct {
	Provider  string `toml:"provider"`
	Model     string `toml:"model"`
	APIKeyEnv string `toml:"api_key_env"`
}

type Scopes struct {
	Predefined []string `toml:"predefined"`
}

type Day struct {
	DataDir string `toml:"data_dir"`
}

type Config struct {
	Ai     Ai     `toml:"ai"`
	Scopes Scopes `toml:"scopes"`
	Day    Day    `toml:"day"`
}

const (
	configPath    = ".config/day/config.toml"
	defaultDBPath = ".local/share/day"
)

func GetConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return &Config{}, fmt.Errorf("could not resolve UserHomeDir: %w", err)
	}

	configPathAbs := filepath.Join(homeDir, configPath)

	c := Config{}
	c.Day.DataDir = defaultDBPath

	if _, err := os.Stat(configPathAbs); os.IsNotExist(err) {
		// silent error here. no config is no error
		return &c, nil
	}

	_, err = toml.DecodeFile(configPathAbs, &c)
	if err != nil {
		return &c, fmt.Errorf("could not decode config file: %w", err)
	}

	return &c, nil
}
