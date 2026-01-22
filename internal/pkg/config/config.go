package config

import (
	"time"
)

// Config defines the configuration interface that any provider must implement.
type Config interface {
	GetString(key string) string
	GetInt(key string) int
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetStringSlice(key string) []string
	GetStringMapString(key string) map[string]string
	GetSection(prefix string) Config
	IsSet(key string) bool
}
