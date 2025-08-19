/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | getters.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"fmt"
	"maps"
	"strconv"
	"strings"
	"time"
)

// GetString retrieves a string value with an optional default.
// Supports both flat keys ("key") and nested keys with dot notation ("server.host").
func (c *Config) GetString(key string, defaultValue ...string) string {
	if c == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return ""
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if value, exists := c.data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}

		return fmt.Sprintf("%v", value)
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		if value, exists := c.getNestedValueUnsafe(key); exists {
			if str, ok := value.(string); ok {
				return str
			}

			return fmt.Sprintf("%v", value)
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// GetInt retrieves an integer value with an optional default.
// Supports both flat keys ("key") and nested keys with dot notation ("server.port").
func (c *Config) GetInt(key string, defaultValue ...int) int {
	if c == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		case float32:
			return int(v)
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				return parsed
			}
		}
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		if value, exists := c.getNestedValueUnsafe(key); exists {
			switch v := value.(type) {
			case int:
				return v
			case int64:
				return int(v)
			case float64:
				return int(v)
			case float32:
				return int(v)
			case string:
				if parsed, err := strconv.Atoi(v); err == nil {
					return parsed
				}
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

// GetFloat64 retrieves a float64 value with an optional default.
// Supports both flat keys ("key") and nested keys with dot notation ("server.load").
func (c *Config) GetFloat64(key string, defaultValue ...float64) float64 {
	if c == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0.0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case string:
			if parsed, err := strconv.ParseFloat(v, 64); err == nil {
				return parsed
			}
		}
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		if value, exists := c.getNestedValueUnsafe(key); exists {
			switch v := value.(type) {
			case float64:
				return v
			case float32:
				return float64(v)
			case int:
				return float64(v)
			case int64:
				return float64(v)
			case string:
				if parsed, err := strconv.ParseFloat(v, 64); err == nil {
					return parsed
				}
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0.0
}

// GetBool retrieves a boolean value with an optional default.
// Supports both flat keys ("key") and nested keys with dot notation ("server.debug").
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	if c == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case bool:
			return v
		case string:
			if parsed, err := strconv.ParseBool(v); err == nil {
				return parsed
			}
		}
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		if value, exists := c.getNestedValueUnsafe(key); exists {
			switch v := value.(type) {
			case bool:
				return v
			case string:
				if parsed, err := strconv.ParseBool(v); err == nil {
					return parsed
				}
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return false
}

// GetDuration retrieves a duration value with an optional default
// It supports string representations like "30s" or "1m"
// and also accepts integers and floats representing seconds.
// Supports both flat keys ("key") and nested keys with dot notation ("server.timeout").
func (c *Config) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	if c == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case string:
			if parsed, err := time.ParseDuration(v); err == nil {
				return parsed
			}

			if seconds, err := strconv.Atoi(v); err == nil {
				return time.Duration(seconds) * time.Second
			}
		case int:
			return time.Duration(v) * time.Second
		case int64:
			return time.Duration(v) * time.Second
		case float64:
			return time.Duration(v * float64(time.Second))
		}
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		if value, exists := c.getNestedValueUnsafe(key); exists {
			switch v := value.(type) {
			case string:
				if parsed, err := time.ParseDuration(v); err == nil {
					return parsed
				}

				if seconds, err := strconv.Atoi(v); err == nil {
					return time.Duration(seconds) * time.Second
				}
			case int:
				return time.Duration(v) * time.Second
			case int64:
				return time.Duration(v) * time.Second
			case float64:
				return time.Duration(v * float64(time.Second))
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

// GetStringSlice retrieves a string slice value.
// Supports both flat keys ("key") and nested keys with dot notation ("server.features").
func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	if c == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}

		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case []string:
			return v
		case []any:
			result := make([]string, len(v))
			for i, item := range v {
				result[i] = fmt.Sprintf("%v", item)
			}

			return result
		case string:
			return strings.Split(v, ",")
		}
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		if value, exists := c.getNestedValueUnsafe(key); exists {
			switch v := value.(type) {
			case []string:
				return v
			case []any:
				result := make([]string, len(v))
				for i, item := range v {
					result[i] = fmt.Sprintf("%v", item)
				}

				return result
			case string:
				return strings.Split(v, ",")
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return []string{}
}

// GetNestedMap returns a nested map at the specified path.
// Returns nil if the path doesn't exist or doesn't point to a map.
func (c *Config) GetNestedMap(key string) map[string]any {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.getNestedValueUnsafe(key); exists {
		if nestedMap, ok := value.(map[string]any); ok {
			// Return a copy to maintain thread safety
			result := make(map[string]any)
			for k, v := range nestedMap {
				result[k] = v
			}

			return result
		}
	}

	return nil
}

// GetNestedKeys returns all keys that start with the given prefix.
// For example, GetNestedKeys("server") might return ["server.host", "server.port", "server.debug"].
func (c *Config) GetNestedKeys(prefix string) []string {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var keys []string

	// Check if prefix exists as a map
	if value, exists := c.data[prefix]; exists {
		if nestedMap, ok := value.(map[string]any); ok {
			for key := range nestedMap {
				keys = append(keys, prefix+"."+key)
			}
		}
	}

	return keys
}

// GetAll returns a copy of all configuration data.
func (c *Config) GetAll() map[string]any {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]any, len(c.data))
	maps.Copy(result, c.data)

	return result
}
