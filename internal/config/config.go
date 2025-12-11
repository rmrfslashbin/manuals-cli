// Package config handles CLI configuration.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the CLI configuration.
type Config struct {
	// APIBaseURL is the base URL of the Manuals API.
	APIBaseURL string `mapstructure:"api_url"`

	// APIKey is the API key for authentication.
	APIKey string `mapstructure:"api_key"`

	// OutputFormat is the default output format (json, table, text).
	OutputFormat string `mapstructure:"output_format"`
}

// Load reads configuration from file and environment.
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("api_url", "http://localhost:8080")
	v.SetDefault("output_format", "table")

	// Config file locations
	v.SetConfigName(".manuals")
	v.SetConfigType("yaml")

	// Look in home directory
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(home)
	}

	// Look in current directory
	v.AddConfigPath(".")

	// Also check XDG config
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		v.AddConfigPath(filepath.Join(xdgConfig, "manuals"))
	} else if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config", "manuals"))
	}

	// Environment variables
	v.SetEnvPrefix("MANUALS")
	v.AutomaticEnv()

	// Bind specific environment variables to config keys
	_ = v.BindEnv("api_url", "MANUALS_API_URL")
	_ = v.BindEnv("api_key", "MANUALS_API_KEY")
	_ = v.BindEnv("output_format", "MANUALS_OUTPUT_FORMAT")

	// Read config file (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %w", err)
	}

	return &cfg, nil
}

// Validate checks that required configuration is present.
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key required: set MANUALS_API_KEY or add api_key to config file")
	}
	return nil
}
