/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | types.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"errors"
	"sync"
)

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
	Format         Format                     // Configuration file format (auto-detected if not specified)
	IgnoreEnv      bool                       // If true, skip environment variable override
	RequiredKeys   []string                   // Keys that must be present after loading
	DefaultValues  map[string]any             // Default values applied before loading file
	ValidationFunc func(map[string]any) error // Custom validation function
}

// Custom errors.
var (
	ErrInvalidFormat      = errors.New("invalid configuration format")
	ErrFileNotFound       = errors.New("configuration file not found")
	ErrInvalidKey         = errors.New("invalid configuration key")
	ErrRequiredKeyMissing = errors.New("required configuration key is missing")
	ErrConfigNil          = errors.New("configuration is nil")
)
