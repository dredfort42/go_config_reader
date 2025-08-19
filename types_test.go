/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | types_test.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormat_Constants tests the Format constants
func TestFormat_Constants(t *testing.T) {
	assert.Equal(t, Format(0), FormatINI)
	assert.Equal(t, Format(1), FormatJSON)
	assert.Equal(t, Format(2), FormatYAML)
}

// TestLoadOptions_Struct tests the LoadOptions struct
func TestLoadOptions_Struct(t *testing.T) {
	opts := &LoadOptions{
		Format:         FormatJSON,
		IgnoreEnv:      true,
		RequiredKeys:   []string{"key1", "key2"},
		DefaultValues:  map[string]any{"default": "value"},
		ValidationFunc: func(data map[string]any) error { return nil },
	}

	assert.Equal(t, FormatJSON, opts.Format)
	assert.True(t, opts.IgnoreEnv)
	assert.Len(t, opts.RequiredKeys, 2)
	assert.Equal(t, "value", opts.DefaultValues["default"])
	assert.NotNil(t, opts.ValidationFunc)
}

// TestErrors_Constants tests the error constants
func TestErrors_Constants(t *testing.T) {
	assert.Equal(t, "invalid configuration format", ErrInvalidFormat.Error())
	assert.Equal(t, "configuration file not found", ErrFileNotFound.Error())
	assert.Equal(t, "invalid configuration key", ErrInvalidKey.Error())
	assert.Equal(t, "required configuration key is missing", ErrRequiredKeyMissing.Error())
	assert.Equal(t, "configuration is nil", ErrConfigNil.Error())
}

// TestConfig_Struct tests the Config struct
func TestConfig_Struct(t *testing.T) {
	c, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.data)
	assert.Equal(t, 0, len(c.data))
}

// TestOption_Function tests the Option function type
func TestOption_Function(t *testing.T) {
	// Test option function that sets initial data
	initOption := func(c *Config) error {
		c.data["initial"] = "value"
		return nil
	}

	c, err := New(initOption)
	assert.NoError(t, err)
	assert.Equal(t, "value", c.GetString("initial"))
}

// TestOption_Function_Error tests the Option function with error
func TestOption_Function_Error(t *testing.T) {
	// Test option function that returns error
	errorOption := func(c *Config) error {
		return ErrInvalidFormat
	}

	c, err := New(errorOption)
	assert.Error(t, err)
	assert.Nil(t, c)
	assert.Contains(t, err.Error(), "failed to apply option")
	assert.Contains(t, err.Error(), ErrInvalidFormat.Error())
}
