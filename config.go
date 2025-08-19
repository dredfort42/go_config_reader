/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | config.go
	::  ::          ::  ::    Created  | 2025-08-07
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// New creates a new Config instance.
func New(options ...Option) (*Config, error) {
	c := &Config{
		data: make(map[string]any),
	}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply option: %w", err)
		}
	}

	return c, nil
}

// LoadFromFile loads configuration from a file with specified options.
// This method is thread-safe and prevents race conditions during loading.
func (c *Config) LoadFromFile(filePath string, opts *LoadOptions) error {
	if opts == nil {
		opts = &LoadOptions{}
	}

	// Check if file exists before acquiring lock
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrFileNotFound, filePath)
	}

	// Determine format from file extension if not specified
	format := opts.Format
	if format == 0 {
		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".json":
			format = FormatJSON
		case ".yaml", ".yml":
			format = FormatYAML
		default:
			format = FormatINI
		}
	}

	// Read file outside of lock to minimize lock time
	// #nosec G304
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse configuration data outside of lock
	var configData map[string]any

	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case FormatINI:
		// Create a temporary config instance for parsing INI
		tempConfig := &Config{}
		configData = tempConfig.parseINI(string(data))
	default:
		return fmt.Errorf("%w: unsupported format", ErrInvalidFormat)
	}

	// Now acquire lock and update configuration atomically
	c.mu.Lock()
	defer c.mu.Unlock()

	// Replace existing data
	c.data = configData

	// Apply default values only for keys that don't exist
	c.applyDefaultsUnsafe(opts.DefaultValues)

	// Override with environment variables unless disabled
	if !opts.IgnoreEnv {
		c.loadFromEnvironmentUnsafe()
	}

	// Validate required keys
	if err := c.validateRequiredKeysUnsafe(opts.RequiredKeys); err != nil {
		return err
	}

	// Apply custom validation
	if opts.ValidationFunc != nil {
		// Create a copy for validation to avoid exposing internal state
		dataCopy := make(map[string]any)
		maps.Copy(dataCopy, c.data)

		if err := opts.ValidationFunc(dataCopy); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// Has checks if a configuration key exists.
// Supports both flat keys ("key") and nested keys with dot notation ("server.host").
func (c *Config) Has(key string) bool {
	if c == nil {
		return false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	// Try flat key first
	if _, exists := c.data[key]; exists {
		return true
	}

	// If flat key doesn't exist and key contains dots, try nested access
	if strings.Contains(key, ".") {
		return c.hasNestedKeyUnsafe(key)
	}

	return false
}

// Keys returns all configuration keys.
func (c *Config) Keys() []string {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}

	return keys
}

// LoadFromMap loads configuration from a map.
func (c *Config) LoadFromMap(data map[string]any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	maps.Copy(c.data, data)
}

// Clear removes all configuration data.
func (c *Config) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[string]any)
}

// Size returns the number of configuration keys.
func (c *Config) Size() int {
	if c == nil {
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

// IsEmpty returns true if the configuration is empty.
func (c *Config) IsEmpty() bool {
	if c == nil {
		return true
	}

	return c.Size() == 0
}

// String returns a string representation of the configuration.
func (c *Config) String() string {
	if c == nil {
		return "Config is nil"
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var sb strings.Builder
	for key, value := range c.data {
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	return sb.String()
}
