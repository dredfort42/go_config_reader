/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | helpers_test.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHelpers_Comprehensive tests helper functions
func TestHelpers_Comprehensive(t *testing.T) {
	t.Run("ApplyDefaults", func(t *testing.T) {
		c, err := New()
		require.NoError(t, err)

		// Set some initial values
		c.Set("existing_key", "existing_value")
		c.Set("nested.existing", "nested_existing")

		// Apply defaults (should not override existing values)
		defaults := map[string]any{
			"existing_key":    "default_value",      // Should not override
			"new_key":         "new_default",        // Should be set
			"nested.existing": "nested_default",     // Should not override
			"nested.new":      "nested_new_default", // Should be set
			"deep.nested.key": "deep_default",       // Should be set
		}

		c.applyDefaultsUnsafe(defaults)

		// Check that existing values were not overridden
		assert.Equal(t, "existing_value", c.GetString("existing_key"))
		assert.Equal(t, "nested_existing", c.GetString("nested.existing"))

		// Check that new defaults were applied
		assert.Equal(t, "new_default", c.GetString("new_key"))
		assert.Equal(t, "nested_new_default", c.GetString("nested.new"))
		assert.Equal(t, "deep_default", c.GetString("deep.nested.key"))
	})

	t.Run("LoadFromEnvironment", func(t *testing.T) {
		c, err := New()
		require.NoError(t, err)

		// Set some initial config values
		c.Set("PORT", "8080")
		c.Set("HOST", "localhost")
		c.Set("DEBUG", "false")
		c.Set("NOT_IN_ENV", "config_value")

		// Set environment variables
		os.Setenv("PORT", "9090")
		os.Setenv("DEBUG", "true")
		// HOST is not set in environment
		// NOT_IN_ENV is not set in environment

		defer func() {
			os.Unsetenv("PORT")
			os.Unsetenv("DEBUG")
		}()

		// Load from environment
		c.loadFromEnvironmentUnsafe()

		// Check that environment values override config values
		assert.Equal(t, "9090", c.GetString("PORT"))
		assert.Equal(t, "true", c.GetString("DEBUG"))

		// Check that values without env vars remain unchanged
		assert.Equal(t, "localhost", c.GetString("HOST"))
		assert.Equal(t, "config_value", c.GetString("NOT_IN_ENV"))
	})

	t.Run("ValidateRequiredKeys", func(t *testing.T) {
		c, err := New()
		require.NoError(t, err)

		// Set up test data
		c.Set("app_name", "test-app")
		c.Set("port", 8080)
		c.Set("database.host", "localhost")
		c.Set("database.port", 5432)

		t.Run("AllRequiredKeysPresent", func(t *testing.T) {
			requiredKeys := []string{"app_name", "port", "database.host", "database.port"}
			err := c.validateRequiredKeysUnsafe(requiredKeys)
			assert.NoError(t, err)
		})

		t.Run("MissingFlatKey", func(t *testing.T) {
			requiredKeys := []string{"app_name", "missing_key"}
			err := c.validateRequiredKeysUnsafe(requiredKeys)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "required configuration key is missing")
			assert.Contains(t, err.Error(), "missing_key")
		})

		t.Run("MissingNestedKey", func(t *testing.T) {
			requiredKeys := []string{"database.host", "database.missing"}
			err := c.validateRequiredKeysUnsafe(requiredKeys)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "required configuration key is missing")
			assert.Contains(t, err.Error(), "database.missing")
		})

		t.Run("EmptyRequiredKeysList", func(t *testing.T) {
			err := c.validateRequiredKeysUnsafe([]string{})
			assert.NoError(t, err)
		})

		t.Run("NilRequiredKeysList", func(t *testing.T) {
			err := c.validateRequiredKeysUnsafe(nil)
			assert.NoError(t, err)
		})

		t.Run("MixedFlatAndNestedKeys", func(t *testing.T) {
			// Test when a flat key exists with same name as nested path
			c.Set("database.connection", "flat_value")
			c.Set("database", map[string]any{
				"connection": map[string]any{
					"pool": "nested_value",
				},
			})

			// Should find the flat key
			requiredKeys := []string{"database.connection"}
			err := c.validateRequiredKeysUnsafe(requiredKeys)
			assert.NoError(t, err)

			// Should not find non-existent nested key
			requiredKeys = []string{"database.connection.missing"}
			err = c.validateRequiredKeysUnsafe(requiredKeys)
			assert.Error(t, err)
		})
	})
}

// TestHelpers_EdgeCases tests edge cases for helper methods
func TestHelpers_EdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test environment variable loading with complex values
	// First add a key to the config, then set the environment variable
	c.Set("TEST_COMPLEX_JSON", "initial_value")

	os.Setenv("TEST_COMPLEX_JSON", `{"nested": "value"}`)
	defer os.Unsetenv("TEST_COMPLEX_JSON")

	c.loadFromEnvironmentUnsafe()
	// The environment loader replaces existing keys with environment values
	assert.Equal(t, `{"nested": "value"}`, c.GetString("TEST_COMPLEX_JSON"))

	// Test applyDefaultsUnsafe with nil defaults
	c.applyDefaultsUnsafe(nil)
	// Should not panic

	// Test validateRequiredKeysUnsafe with empty keys
	err = c.validateRequiredKeysUnsafe([]string{})
	assert.NoError(t, err)

	// Test validateRequiredKeysUnsafe with nil keys
	err = c.validateRequiredKeysUnsafe(nil)
	assert.NoError(t, err)
}

// Benchmark Tests for helpers.go functions

func BenchmarkConfig_ApplyDefaultsUnsafe(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	defaults := map[string]any{
		"key1":        "value1",
		"key2":        42,
		"key3":        true,
		"nested.key":  "nested_value",
		"deep.nested": "deep_value",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.applyDefaultsUnsafe(defaults)
	}
}

func BenchmarkConfig_LoadFromEnvironmentUnsafe(b *testing.B) {
	// Set up some environment variables
	os.Setenv("BENCH_TEST_VAR1", "value1")
	os.Setenv("BENCH_TEST_VAR2", "42")
	os.Setenv("BENCH_TEST_VAR3", "true")
	defer func() {
		os.Unsetenv("BENCH_TEST_VAR1")
		os.Unsetenv("BENCH_TEST_VAR2")
		os.Unsetenv("BENCH_TEST_VAR3")
	}()

	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.loadFromEnvironmentUnsafe()
	}
}

func BenchmarkConfig_ValidateRequiredKeysUnsafe(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"required1": "value1",
		"required2": "value2",
		"nested": map[string]any{
			"required3": "value3",
		},
	})

	requiredKeys := []string{"required1", "required2", "nested.required3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := c.validateRequiredKeysUnsafe(requiredKeys)
		if err != nil {
			b.Fatal(err)
		}
	}
}
