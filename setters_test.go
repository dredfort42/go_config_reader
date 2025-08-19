/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | setters_test.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSetNestedDefaults_Comprehensive tests SetNestedDefaults method
func TestSetNestedDefaults_Comprehensive(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	t.Run("SetDefaultsWithoutOverriding", func(t *testing.T) {
		// Set some existing values
		c.Set("app.name", "existing-app")
		c.Set("app.version", "1.0.0")
		c.Set("database.host", "existing-host")

		// Set nested defaults
		defaults := map[string]any{
			"app.name":      "default-app",  // Should not override
			"app.debug":     true,           // Should be set
			"app.port":      8080,           // Should be set
			"database.host": "default-host", // Should not override
			"database.port": 5432,           // Should be set
			"cache.enabled": true,           // Should be set
			"cache.ttl":     300,            // Should be set
			"flat_key":      "flat_value",   // Should be set
		}

		c.SetNestedDefaults(defaults)

		// Check existing values weren't overridden
		assert.Equal(t, "existing-app", c.GetString("app.name"))
		assert.Equal(t, "1.0.0", c.GetString("app.version"))
		assert.Equal(t, "existing-host", c.GetString("database.host"))

		// Check new defaults were set
		assert.True(t, c.GetBool("app.debug"))
		assert.Equal(t, 8080, c.GetInt("app.port"))
		assert.Equal(t, 5432, c.GetInt("database.port"))
		assert.True(t, c.GetBool("cache.enabled"))
		assert.Equal(t, 300, c.GetInt("cache.ttl"))
		assert.Equal(t, "flat_value", c.GetString("flat_key"))
	})

	t.Run("EmptyDefaults", func(t *testing.T) {
		c.Clear()
		c.Set("existing", "value")

		// Apply empty defaults
		c.SetNestedDefaults(map[string]any{})

		// Nothing should change
		assert.Equal(t, "value", c.GetString("existing"))
		assert.Equal(t, 1, c.Size())
	})

	t.Run("NilDefaults", func(t *testing.T) {
		c.Clear()
		c.Set("existing", "value")

		// Apply nil defaults
		c.SetNestedDefaults(nil)

		// Nothing should change
		assert.Equal(t, "value", c.GetString("existing"))
		assert.Equal(t, 1, c.Size())
	})

	t.Run("DeepNestedDefaults", func(t *testing.T) {
		c.Clear()

		// Set deep nested defaults
		defaults := map[string]any{
			"level1.level2.level3.level4": "deep_value",
			"level1.level2.other":         "other_value",
			"level1.different":            "different_value",
		}

		c.SetNestedDefaults(defaults)

		// Check all levels were created
		assert.Equal(t, "deep_value", c.GetString("level1.level2.level3.level4"))
		assert.Equal(t, "other_value", c.GetString("level1.level2.other"))
		assert.Equal(t, "different_value", c.GetString("level1.different"))

		// Check intermediate maps exist
		assert.NotNil(t, c.GetNestedMap("level1"))
		assert.NotNil(t, c.GetNestedMap("level1.level2"))
		assert.NotNil(t, c.GetNestedMap("level1.level2.level3"))
	})
}

// TestSetters_EdgeCases tests edge cases for setter methods
func TestSetters_EdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test Set with nil value
	c.Set("nil_key", nil)
	assert.True(t, c.Has("nil_key"))

	// Test Set with empty key
	c.Set("", "empty_key_value")
	assert.Equal(t, "empty_key_value", c.GetString(""))

	// Test SetNestedDefaults with nil config (shouldn't panic)
	var nilConfig *Config
	assert.NotPanics(t, func() {
		nilConfig.SetNestedDefaults(map[string]any{"key": "value"})
	})

	// Test Set with nil config (shouldn't panic)
	assert.NotPanics(t, func() {
		nilConfig.Set("key", "value")
	})
}

// Benchmark Tests for setters.go functions

func BenchmarkConfig_Set(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("test_key", "test_value")
		c.Set("nested.key", "nested_value")
	}
}

func BenchmarkConfig_SetNestedDefaults(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	defaults := map[string]any{
		"key1":               "value1",
		"nested.key2":        "value2",
		"deep.nested.key3":   "value3",
		"another.branch.key": "value4",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.SetNestedDefaults(defaults)
	}
}
