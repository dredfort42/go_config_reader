/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | compatibility.go
	::  ::          ::  ::    Created  | 2025-08-07
		  ::::  ::::          Modified | 2025-08-07

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"fmt"
	"sync"
	"time"
)

// Global instance for backward compatibility and convenience.
var globalConfig *Config
var globalConfigOnce sync.Once

// initGlobalConfig initializes the global config instance.
func initGlobalConfig() {
	globalConfigOnce.Do(func() {
		var err error

		globalConfig, err = New()
		if err != nil {
			panic(fmt.Sprintf("failed to initialize global config: %v", err))
		}
	})
}

// Global API functions for convenience and backward compatibility

// Set sets a configuration value in the global config instance.
func Set(key string, value any) {
	initGlobalConfig()
	globalConfig.Set(key, value)
}

// GetString retrieves a string value from the global config instance.
func GetString(key string, defaultValue ...string) string {
	initGlobalConfig()

	return globalConfig.GetString(key, defaultValue...)
}

// GetInt retrieves an integer value from the global config instance.
func GetInt(key string, defaultValue ...int) int {
	initGlobalConfig()

	return globalConfig.GetInt(key, defaultValue...)
}

// GetFloat64 retrieves a float64 value from the global config instance.
func GetFloat64(key string, defaultValue ...float64) float64 {
	initGlobalConfig()

	return globalConfig.GetFloat64(key, defaultValue...)
}

// GetBool retrieves a boolean value from the global config instance.
func GetBool(key string, defaultValue ...bool) bool {
	initGlobalConfig()

	return globalConfig.GetBool(key, defaultValue...)
}

// GetDuration retrieves a duration value from the global config instance.
func GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	initGlobalConfig()

	return globalConfig.GetDuration(key, defaultValue...)
}

// GetStringSlice retrieves a string slice value from the global config instance.
func GetStringSlice(key string, defaultValue ...[]string) []string {
	initGlobalConfig()

	return globalConfig.GetStringSlice(key, defaultValue...)
}

// Has checks if a configuration key exists in the global config instance.
func Has(key string) bool {
	initGlobalConfig()

	return globalConfig.Has(key)
}

// Keys returns all configuration keys from the global config instance.
func Keys() []string {
	initGlobalConfig()

	return globalConfig.Keys()
}

// Load loads configuration from a file into the global config instance.
func Load(filePath string, opts *LoadOptions) error {
	initGlobalConfig()

	return globalConfig.LoadFromFile(filePath, opts)
}

// MustLoad loads configuration from a file into a new Config instance and panics on error.
func MustLoad(filePath string, opts *LoadOptions) *Config {
	c, err := New()
	if err != nil {
		panic(fmt.Sprintf("failed to create config: %v", err))
	}

	if err := c.LoadFromFile(filePath, opts); err != nil {
		panic(fmt.Sprintf("failed to load config from %s: %v", filePath, err))
	}

	return c
}

// LoadWithDefaults loads configuration from a file with default values into a new Config instance.
func LoadWithDefaults(filePath string, defaults map[string]any) *Config {
	opts := &LoadOptions{
		DefaultValues: defaults,
	}

	return MustLoad(filePath, opts)
}

// String returns a string representation of the global config instance.
func String() string {
	initGlobalConfig()

	return globalConfig.String()
}
