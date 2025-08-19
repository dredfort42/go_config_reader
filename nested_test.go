/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | nested_test.go
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

// TestNestedOperations_Comprehensive tests nested key operations
func TestNestedOperations_Comprehensive(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	t.Run("SetAndGetNestedValues", func(t *testing.T) {
		// Test setting nested values
		c.Set("level1.level2.level3", "deep_value")
		c.Set("level1.level2.another", "another_value")
		c.Set("level1.different", "different_value")

		// Test getting nested values
		assert.Equal(t, "deep_value", c.GetString("level1.level2.level3"))
		assert.Equal(t, "another_value", c.GetString("level1.level2.another"))
		assert.Equal(t, "different_value", c.GetString("level1.different"))

		// Test Has with nested keys
		assert.True(t, c.Has("level1.level2.level3"))
		assert.True(t, c.Has("level1.level2.another"))
		assert.True(t, c.Has("level1.different"))
		assert.False(t, c.Has("level1.level2.nonexistent"))
		assert.False(t, c.Has("nonexistent.key"))
	})

	t.Run("OverwriteNonMapValues", func(t *testing.T) {
		// Set a simple value first
		c.Set("config.simple", "simple_value")
		assert.Equal(t, "simple_value", c.GetString("config.simple"))

		// Now set a nested value that requires creating a map where simple value was
		c.Set("config.simple.nested", "nested_value")
		assert.Equal(t, "nested_value", c.GetString("config.simple.nested"))

		// The original simple value should be replaced by a map
		// Getting the parent should return empty string since it's now a map
		nestedMap := c.GetNestedMap("config.simple")
		require.NotNil(t, nestedMap)
		assert.Equal(t, "nested_value", nestedMap["nested"])
	})

	t.Run("CreateDeepNesting", func(t *testing.T) {
		// Create deep nesting in one operation
		c.Set("app.database.connection.pool.size", 10)
		c.Set("app.database.connection.pool.timeout", 30)
		c.Set("app.database.connection.host", "localhost")
		c.Set("app.cache.redis.host", "redis-host")

		// Verify all levels can be accessed
		assert.Equal(t, 10, c.GetInt("app.database.connection.pool.size"))
		assert.Equal(t, 30, c.GetInt("app.database.connection.pool.timeout"))
		assert.Equal(t, "localhost", c.GetString("app.database.connection.host"))
		assert.Equal(t, "redis-host", c.GetString("app.cache.redis.host"))

		// Verify intermediate maps exist
		poolMap := c.GetNestedMap("app.database.connection.pool")
		require.NotNil(t, poolMap)
		assert.Equal(t, 10, poolMap["size"])
		assert.Equal(t, 30, poolMap["timeout"])

		connectionMap := c.GetNestedMap("app.database.connection")
		require.NotNil(t, connectionMap)
		assert.Equal(t, "localhost", connectionMap["host"])

		databaseMap := c.GetNestedMap("app.database")
		require.NotNil(t, databaseMap)
		assert.NotNil(t, databaseMap["connection"])
	})

	t.Run("GetNestedKeys", func(t *testing.T) {
		// Clear and set up test data
		c.Clear()
		c.Set("server.host", "localhost")
		c.Set("server.port", 8080)
		c.Set("server.ssl", true)
		c.Set("database.host", "db-host")
		c.Set("other", "value")

		// Test getting nested keys for a section
		serverKeys := c.GetNestedKeys("server")
		assert.Contains(t, serverKeys, "server.host")
		assert.Contains(t, serverKeys, "server.port")
		assert.Contains(t, serverKeys, "server.ssl")
		assert.Len(t, serverKeys, 3)

		// Test getting keys for section with only one key
		databaseKeys := c.GetNestedKeys("database")
		assert.Contains(t, databaseKeys, "database.host")
		assert.Len(t, databaseKeys, 1)

		// Test getting keys for non-existent section
		nonExistentKeys := c.GetNestedKeys("nonexistent")
		assert.Empty(t, nonExistentKeys)

		// Test getting keys for a non-map value
		otherKeys := c.GetNestedKeys("other")
		assert.Empty(t, otherKeys)
	})

	t.Run("EmptyAndSpecialKeys", func(t *testing.T) {
		c.Clear()

		// Test empty key components - these create nested structures
		c.Set(".", "dot_value")
		c.Set(".key", "leading_dot")
		c.Set("key.", "trailing_dot")
		c.Set("key..double", "double_dot")

		assert.Equal(t, "dot_value", c.GetString("."))
		assert.Equal(t, "leading_dot", c.GetString(".key"))

		// key. creates a nested structure where "" is the key
		// key..double creates key -> "" -> "double" = "double_dot"
		// So getting "key." gets the nested map, not the string
		assert.True(t, c.Has("key."))
		assert.Equal(t, "double_dot", c.GetString("key..double"))
	})

	t.Run("PathTraversalSafety", func(t *testing.T) {
		c.Clear()

		// Set nested value normally
		c.Set("safe.nested.value", "safe_content")

		// Try to access through a value that's not a map
		c.Set("not_map", "string_value")

		// These should return default values since path can't be traversed
		assert.Equal(t, "", c.GetString("not_map.fake.path"))
		assert.Equal(t, 0, c.GetInt("not_map.fake.path"))
		assert.False(t, c.GetBool("not_map.fake.path"))
		assert.False(t, c.Has("not_map.fake.path"))
		assert.Nil(t, c.GetNestedMap("not_map.fake"))
	})
}

// TestNested_EdgeCases tests edge cases for nested operations
func TestNested_EdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with nil config
	var nilConfig *Config
	_, exists := nilConfig.getNestedValueUnsafe("any.key")
	assert.False(t, exists)

	exists = nilConfig.hasNestedKeyUnsafe("any.key")
	assert.False(t, exists)

	// Test with empty string key
	_, exists = c.getNestedValueUnsafe("")
	assert.False(t, exists)

	exists = c.hasNestedKeyUnsafe("")
	assert.False(t, exists)

	// Test with single dot
	c.Set(".", "dot_value")
	value, exists := c.getNestedValueUnsafe(".")
	assert.True(t, exists)
	assert.Equal(t, "dot_value", value)
}

// Benchmark Tests for nested.go functions

func BenchmarkConfig_GetNestedValueUnsafe(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	// Create deep nested structure
	c.LoadFromMap(map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": map[string]any{
					"level4": map[string]any{
						"level5": "deep_value",
					},
				},
			},
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.getNestedValueUnsafe("level1.level2.level3.level4.level5")
		c.getNestedValueUnsafe("level1.level2")
		c.getNestedValueUnsafe("nonexistent.path")
	}
}

func BenchmarkConfig_HasNestedKeyUnsafe(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	// Create nested structure
	c.LoadFromMap(map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3":  "value",
				"another": "value2",
			},
			"sibling": "value3",
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.hasNestedKeyUnsafe("level1.level2.level3")
		c.hasNestedKeyUnsafe("level1.sibling")
		c.hasNestedKeyUnsafe("nonexistent.path")
	}
}
