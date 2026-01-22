package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type viperConfig struct {
	v *viper.Viper
}

// Load loads configuration from file and environment variables using Viper provider.
func Load() (*viperConfig, error) {
	v := viper.New()
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("..")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &viperConfig{v: v}, nil
}

func (c *viperConfig) GetString(key string) string {
	return c.v.GetString(key)
}

func (c *viperConfig) GetInt(key string) int {
	return c.v.GetInt(key)
}

func (c *viperConfig) GetInt64(key string) int64 {
	return c.v.GetInt64(key)
}

func (c *viperConfig) GetFloat64(key string) float64 {
	return c.v.GetFloat64(key)
}

func (c *viperConfig) GetBool(key string) bool {
	return c.v.GetBool(key)
}

func (c *viperConfig) GetDuration(key string) time.Duration {
	return c.v.GetDuration(key)
}

func (c *viperConfig) GetStringSlice(key string) []string {
	return c.v.GetStringSlice(key)
}

func (c *viperConfig) GetStringMapString(key string) map[string]string {
	result := make(map[string]string)
	for k, v := range c.v.GetStringMap(key) {
		if str, ok := v.(string); ok {
			result[k] = str
		}
	}
	return result
}

func (c *viperConfig) GetSection(prefix string) Config {
	return &viperConfig{v: c.v.Sub(prefix)}
}

func (c *viperConfig) IsSet(key string) bool {
	return c.v.IsSet(key)
}
