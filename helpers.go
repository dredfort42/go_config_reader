/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | helpers.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"fmt"
	"os"
	"strings"
)

// applyDefaultsUnsafe applies default values for keys that don't exist.
// This method assumes the caller holds the write lock.
func (c *Config) applyDefaultsUnsafe(defaults map[string]any) {
	if c == nil || defaults == nil {
		return
	}

	for key, value := range defaults {
		if strings.Contains(key, ".") {
			// Nested key - check if it exists before setting
			if !c.hasNestedKeyUnsafe(key) {
				c.setNestedValueUnsafe(key, value)
			}
		} else {
			// Flat key - check if it exists before setting
			if _, exists := c.data[key]; !exists {
				c.data[key] = value
			}
		}
	}
}

// loadFromEnvironmentUnsafe loads configuration from environment variables.
// This method assumes the caller holds the write lock.
func (c *Config) loadFromEnvironmentUnsafe() {
	if c == nil {
		return
	}

	for key := range c.data {
		if envValue := os.Getenv(key); envValue != "" {
			c.data[key] = envValue
		}
	}
}

// validateRequiredKeysUnsafe checks if all required keys are present.
// This method assumes the caller holds the appropriate lock.
func (c *Config) validateRequiredKeysUnsafe(requiredKeys []string) error {
	if c == nil {
		return ErrConfigNil
	}

	for _, key := range requiredKeys {
		// Check flat key first
		if _, exists := c.data[key]; exists {
			continue
		}

		// If flat key doesn't exist and key contains dots, try nested access
		if strings.Contains(key, ".") {
			if c.hasNestedKeyUnsafe(key) {
				continue
			}
		}

		// Key not found
		return fmt.Errorf("%w: %q", ErrRequiredKeyMissing, key)
	}

	return nil
}
