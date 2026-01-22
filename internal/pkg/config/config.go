package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config holds the application configuration using Viper.
type Config struct {
	v *viper.Viper
}

// Load loads configuration from file and environment variables.
func Load() (*Config, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Try to read .env file if it exists, but don't fail if it doesn't
	// In Docker mode, environment variables are passed via docker-compose.yaml
	if err := v.ReadInConfig(); err != nil {
		// Only fail if the error is NOT "config file not found"
		// If config file is not found, we'll rely on environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	return &Config{v: v}, nil
}

func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *Config) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *Config) GetStringMapString(key string) map[string]string {
	result := make(map[string]string)
	for k, v := range c.v.GetStringMap(key) {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}

func (c *Config) GetSection(prefix string) *Config {
	return &Config{v: c.v.Sub(prefix)}
}

func (c *Config) IsSet(key string) bool {
	return c.v.IsSet(key)
}
