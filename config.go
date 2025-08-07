/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | config.go
	::  ::          ::  ::    Created  | 2025-08-07
		  ::::  ::::          Modified | 2025-08-07

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

// Package config provides a modern, flexible, and thread-safe configuration library
// for Go applications that supports multiple formats and provides a clean API.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

// Map represents a map of configuration keys and values (backward compatibility).
type Map map[string]string

// File is a global variable that holds the configuration loaded from a file (backward compatibility).
var File Map

// Config represents a configuration manager that provides thread-safe access to configuration values.
type Config struct {
	mu   sync.RWMutex
	data map[string]any
}

// Format represents supported configuration file formats.
type Format int

// Supported configuration formats.
const (
	FormatINI Format = iota
	FormatJSON
	FormatYAML
)

// Option represents a functional option for configuration.
type Option func(*Config) error

// LoadOptions holds options for loading configuration.
type LoadOptions struct {
	Format         Format
	IgnoreEnv      bool
	RequiredKeys   []string
	DefaultValues  map[string]any
	ValidationFunc func(map[string]any) error
}

// Custom errors.
var (
	// Legacy errors for backward compatibility.
	ErrConfigFileNotFound = errors.New("use --config flag to specify the path to configuration file")
	ErrConfigFileClose    = errors.New("could not close configuration file")
	ErrFailedToReadConfig = errors.New("could not read configuration from *.ini file")

	// New specific errors.
	ErrInvalidFormat      = errors.New("invalid configuration format")
	ErrFileNotFound       = errors.New("configuration file not found")
	ErrInvalidKey         = errors.New("invalid configuration key")
	ErrRequiredKeyMissing = errors.New("required configuration key is missing")
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
func (c *Config) LoadFromFile(filePath string, opts *LoadOptions) error {
	if opts == nil {
		opts = &LoadOptions{}
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrFileNotFound, filePath)
	}

	// Determine format from file extension if not specified
	if opts.Format == 0 {
		ext := strings.ToLower(filepath.Ext(filePath))
		switch ext {
		case ".json":
			opts.Format = FormatJSON
		case ".yaml", ".yml":
			opts.Format = FormatYAML
		default:
			opts.Format = FormatINI
		}
	}

	// #nosec G304
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var configData map[string]any

	switch opts.Format {
	case FormatJSON:
		if err := json.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse JSON config: %w", err)
		}
	case FormatYAML:
		if err := yaml.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case FormatINI:
		configData = c.parseINI(string(data))
	default:
		return fmt.Errorf("%w: unsupported format", ErrInvalidFormat)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Apply default values first
	if opts.DefaultValues != nil {
		maps.Copy(c.data, opts.DefaultValues)
	}

	// Load file configuration
	maps.Copy(c.data, configData)

	// Override with environment variables unless disabled
	if !opts.IgnoreEnv {
		c.loadFromEnvironment()
	}

	// Validate required keys
	if err := c.validateRequiredKeys(opts.RequiredKeys); err != nil {
		return err
	}

	// Apply custom validation
	if opts.ValidationFunc != nil {
		if err := opts.ValidationFunc(c.data); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// GetString retrieves a string value with an optional default.
func (c *Config) GetString(key string, defaultValue ...string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.data[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
		// Try to convert to string
		return fmt.Sprintf("%v", value)
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

// GetInt retrieves an integer value with an optional default.
func (c *Config) GetInt(key string, defaultValue ...int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		case string:
			if parsed, err := strconv.Atoi(v); err == nil {
				return parsed
			}
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

// GetFloat64 retrieves a float64 value with an optional default.
func (c *Config) GetFloat64(key string, defaultValue ...float64) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

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

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0.0
}

// GetBool retrieves a boolean value with an optional default.
func (c *Config) GetBool(key string, defaultValue ...bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

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

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return false
}

// GetDuration retrieves a duration value with an optional default
// It supports string representations like "30s" or "1m"
// and also accepts integers and floats representing seconds.
func (c *Config) GetDuration(key string, defaultValue ...time.Duration) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.data[key]; exists {
		switch v := value.(type) {
		case string:
			if parsed, err := time.ParseDuration(v); err == nil {
				return parsed
			}
			// Try parsing as seconds (backward compatibility)
			if seconds, err := strconv.Atoi(v); err == nil {
				return time.Duration(seconds) * time.Second
			}
		case int:
			return time.Duration(v) * time.Second
		case int64:
			return time.Duration(v) * time.Second
		case float64:
			return time.Duration(v) * time.Second
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return 0
}

// GetStringSlice retrieves a string slice value.
func (c *Config) GetStringSlice(key string, defaultValue ...[]string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

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
			// Try to parse comma-separated values
			return strings.Split(v, ",")
		}
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return []string{}
}

// Set sets a configuration value (useful for runtime configuration changes).
func (c *Config) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value
}

// Has checks if a configuration key exists.
func (c *Config) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.data[key]

	return exists
}

// Keys returns all configuration keys.
func (c *Config) Keys() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.data))
	for key := range c.data {
		keys = append(keys, key)
	}

	return keys
}

// GetAll returns a copy of all configuration data.
func (c *Config) GetAll() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]any, len(c.data))
	maps.Copy(result, c.data)

	return result
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
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

// IsEmpty returns true if the configuration is empty.
func (c *Config) IsEmpty() bool {
	return c.Size() == 0
}

// loadFromEnvironment loads configuration from environment variables.
func (c *Config) loadFromEnvironment() {
	for key := range c.data {
		if envValue := os.Getenv(key); envValue != "" {
			c.data[key] = envValue
		}
	}
}

// validateRequiredKeys checks if all required keys are present.
func (c *Config) validateRequiredKeys(requiredKeys []string) error {
	for _, key := range requiredKeys {
		if _, exists := c.data[key]; !exists {
			return fmt.Errorf("%w: %q", ErrRequiredKeyMissing, key)
		}
	}

	return nil
}

// parseINI parses INI format content (backward compatibility).
func (c *Config) parseINI(content string) map[string]any {
	result := make(map[string]any)
	lines := strings.SplitSeq(content, "\n")

	for line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Remove inline comments
		if idx := strings.IndexAny(line, "#;"); idx != -1 {
			line = line[:idx]
			line = strings.TrimSpace(line)
		}

		key, value, found := strings.Cut(line, "=")
		if !found || key == "" {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// Remove surrounding quotes
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		result[key] = value
	}

	return result
}

// String returns a string representation of the configuration.
func (c *Config) String() string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var sb strings.Builder
	for key, value := range c.data {
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, value))
	}

	return sb.String()
}
