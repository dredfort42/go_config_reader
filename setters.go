/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | setters.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"strings"
)

// Set sets a configuration value (useful for runtime configuration changes).
// Supports both flat keys ("key") and nested keys with dot notation ("server.host").
func (c *Config) Set(key string, value any) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// If key contains dots, use nested setting
	if strings.Contains(key, ".") {
		c.setNestedValueUnsafe(key, value)
	} else {
		c.data[key] = value
	}
}

// This is useful for setting up complex default configurations.
func (c *Config) SetNestedDefaults(defaults map[string]any) {
	if c == nil || defaults == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for key, value := range defaults {
		if strings.Contains(key, ".") {
			// Check if the nested key already exists
			if !c.hasNestedKeyUnsafe(key) {
				c.setNestedValueUnsafe(key, value)
			}
		} else {
			// Check if the flat key already exists
			if _, exists := c.data[key]; !exists {
				c.data[key] = value
			}
		}
	}
}

// setNestedValueUnsafe sets a nested value using dot notation (e.g., "server.host").
// This method assumes the caller holds the appropriate lock.
func (c *Config) setNestedValueUnsafe(key string, value any) {
	if c == nil {
		return
	}

	parts := strings.Split(key, ".")
	current := c.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - set the value
			current[part] = value

			return
		}

		// Navigate or create nested structure
		if existing, exists := current[part]; exists {
			if nestedMap, ok := existing.(map[string]any); ok {
				current = nestedMap
			} else {
				// Existing value is not a map, replace it with a map
				current[part] = make(map[string]any)
				if m, ok := current[part].(map[string]any); ok {
					current = m
				}
			}
		} else {
			// Create new nested map
			current[part] = make(map[string]any)
			if m, ok := current[part].(map[string]any); ok {
				current = m
			}
		}
	}
}
