/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | nested.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import "strings"

// getNestedValueUnsafe retrieves a nested value using dot notation (e.g., "server.host").
// This method assumes the caller holds the appropriate lock.
func (c *Config) getNestedValueUnsafe(key string) (any, bool) {
	if c == nil {
		return nil, false
	}

	parts := strings.Split(key, ".")
	current := c.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - get the final value
			if value, exists := current[part]; exists {
				return value, true
			}

			return nil, false
		}

		// Navigate deeper into nested structure
		if value, exists := current[part]; exists {
			if nestedMap, ok := value.(map[string]any); ok {
				current = nestedMap
			} else {
				// Path doesn't lead to a map, can't continue
				return nil, false
			}
		} else {
			// Key doesn't exist in path
			return nil, false
		}
	}

	return nil, false
}

// hasNestedKeyUnsafe checks if a nested key exists using dot notation.
// This method assumes the caller holds the appropriate lock.
func (c *Config) hasNestedKeyUnsafe(key string) bool {
	if c == nil {
		return false
	}

	_, exists := c.getNestedValueUnsafe(key)

	return exists
}
