/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | getters_test.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetters_Comprehensive tests all getter methods with extensive coverage
func TestGetters_Comprehensive(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Setup test data with various types
	testData := map[string]any{
		"string_value":      "hello",
		"int_value":         42,
		"int64_value":       int64(9223372036854775807),
		"float64_value":     3.14159,
		"float32_value":     float32(2.71),
		"bool_true":         true,
		"bool_false":        false,
		"string_int":        "123",
		"string_float":      "45.67",
		"string_bool_true":  "true",
		"string_bool_false": "false",
		"string_bool_yes":   "yes",
		"string_bool_no":    "no",
		"string_bool_on":    "on",
		"string_bool_off":   "off",
		"string_bool_1":     "1",
		"string_bool_0":     "0",
		"duration_string":   "30s",
		"duration_minutes":  "5m",
		"duration_hours":    "2h",
		"duration_complex":  "1h30m45s",
		"duration_int":      30,
		"duration_int64":    int64(60),
		"duration_float":    90.5,
		"string_slice":      []string{"a", "b", "c"},
		"any_slice":         []any{"x", "y", "z"},
		"comma_separated":   "red,green,blue",
		"empty_string":      "",
		"zero_int":          0,
		"zero_float":        0.0,
		"nested": map[string]any{
			"level1": map[string]any{
				"string_val": "nested_string",
				"int_val":    100,
				"bool_val":   true,
			},
		},
	}

	c.LoadFromMap(testData)

	t.Run("GetString", func(t *testing.T) {
		// Test string value
		assert.Equal(t, "hello", c.GetString("string_value"))

		// Test conversion from other types
		assert.Equal(t, "42", c.GetString("int_value"))
		assert.Equal(t, "3.14159", c.GetString("float64_value"))
		assert.Equal(t, "true", c.GetString("bool_true"))

		// Test nested access
		assert.Equal(t, "nested_string", c.GetString("nested.level1.string_val"))

		// Test default value
		assert.Equal(t, "default", c.GetString("nonexistent", "default"))
		assert.Equal(t, "", c.GetString("nonexistent"))

		// Test empty string
		assert.Equal(t, "", c.GetString("empty_string"))
	})

	t.Run("GetInt", func(t *testing.T) {
		// Test int value
		assert.Equal(t, 42, c.GetInt("int_value"))

		// Test conversion from int64
		assert.Equal(t, int(9223372036854775807), c.GetInt("int64_value"))

		// Test conversion from float64
		assert.Equal(t, 3, c.GetInt("float64_value"))
		assert.Equal(t, 2, c.GetInt("float32_value")) // float32(2.71) converts to int 2

		// Test conversion from string
		assert.Equal(t, 123, c.GetInt("string_int"))

		// Test nested access
		assert.Equal(t, 100, c.GetInt("nested.level1.int_val"))

		// Test default value
		assert.Equal(t, 999, c.GetInt("nonexistent", 999))
		assert.Equal(t, 0, c.GetInt("nonexistent"))

		// Test invalid string conversion
		assert.Equal(t, 0, c.GetInt("string_value"))

		// Test zero value
		assert.Equal(t, 0, c.GetInt("zero_int"))
	})

	t.Run("GetFloat64", func(t *testing.T) {
		// Test float64 value
		assert.Equal(t, 3.14159, c.GetFloat64("float64_value"))

		// Test conversion from float32
		assert.InDelta(t, 2.71, c.GetFloat64("float32_value"), 0.01) // Allow more tolerance for float32 conversion

		// Test conversion from int
		assert.Equal(t, 42.0, c.GetFloat64("int_value"))

		// Test conversion from int64
		assert.Equal(t, float64(9223372036854775807), c.GetFloat64("int64_value"))

		// Test conversion from string
		assert.Equal(t, 45.67, c.GetFloat64("string_float"))

		// Test default value
		assert.Equal(t, 99.9, c.GetFloat64("nonexistent", 99.9))
		assert.Equal(t, 0.0, c.GetFloat64("nonexistent"))

		// Test invalid string conversion
		assert.Equal(t, 0.0, c.GetFloat64("string_value"))

		// Test zero value
		assert.Equal(t, 0.0, c.GetFloat64("zero_float"))
	})

	t.Run("GetBool", func(t *testing.T) {
		// Test bool values
		assert.True(t, c.GetBool("bool_true"))
		assert.False(t, c.GetBool("bool_false"))

		// Test string conversions - true cases
		assert.True(t, c.GetBool("string_bool_true"))
		assert.False(t, c.GetBool("string_bool_yes")) // "yes" is not recognized by strconv.ParseBool
		assert.False(t, c.GetBool("string_bool_on"))  // "on" is not recognized by strconv.ParseBool
		assert.True(t, c.GetBool("string_bool_1"))

		// Test string conversions - false cases
		assert.False(t, c.GetBool("string_bool_false"))
		// "no" is not recognized by strconv.ParseBool (returns false as default)
		assert.False(t, c.GetBool("string_bool_no"))
		// "off" is not recognized by strconv.ParseBool (returns false as default)
		assert.False(t, c.GetBool("string_bool_off"))
		assert.False(t, c.GetBool("string_bool_0"))

		// Test nested access
		assert.True(t, c.GetBool("nested.level1.bool_val"))

		// Test default value
		assert.True(t, c.GetBool("nonexistent", true))
		assert.False(t, c.GetBool("nonexistent"))

		// Test invalid string conversion
		assert.False(t, c.GetBool("string_value"))
	})

	t.Run("GetDuration", func(t *testing.T) {
		// Test duration string parsing
		assert.Equal(t, 30*time.Second, c.GetDuration("duration_string"))
		assert.Equal(t, 5*time.Minute, c.GetDuration("duration_minutes"))
		assert.Equal(t, 2*time.Hour, c.GetDuration("duration_hours"))

		expected := 1*time.Hour + 30*time.Minute + 45*time.Second
		assert.Equal(t, expected, c.GetDuration("duration_complex"))

		// Test conversion from int (seconds)
		assert.Equal(t, 30*time.Second, c.GetDuration("duration_int"))

		// Test conversion from int64 (seconds)
		assert.Equal(t, 60*time.Second, c.GetDuration("duration_int64"))

		// Test conversion from float64 (seconds)
		expectedDuration := time.Duration(90.5 * float64(time.Second))
		actual := c.GetDuration("duration_float")
		assert.InDelta(t, float64(expectedDuration), float64(actual), float64(time.Millisecond))

		// Test string number parsing as seconds
		c.Set("duration_string_number", "120")
		assert.Equal(t, 120*time.Second, c.GetDuration("duration_string_number"))

		// Test default value
		defaultDuration := 10 * time.Minute
		assert.Equal(t, defaultDuration, c.GetDuration("nonexistent", defaultDuration))
		assert.Equal(t, time.Duration(0), c.GetDuration("nonexistent"))

		// Test invalid string
		c.Set("invalid_duration", "invalid")
		assert.Equal(t, time.Duration(0), c.GetDuration("invalid_duration"))
	})

	t.Run("GetStringSlice", func(t *testing.T) {
		// Test string slice value
		assert.Equal(t, []string{"a", "b", "c"}, c.GetStringSlice("string_slice"))

		// Test any slice conversion
		assert.Equal(t, []string{"x", "y", "z"}, c.GetStringSlice("any_slice"))

		// Test comma-separated string
		assert.Equal(t, []string{"red", "green", "blue"}, c.GetStringSlice("comma_separated"))

		// Test single string (comma-separated with one item)
		c.Set("single_item", "single")
		assert.Equal(t, []string{"single"}, c.GetStringSlice("single_item"))

		// Test empty string
		assert.Equal(t, []string{""}, c.GetStringSlice("empty_string"))

		// Test default value
		defaultSlice := []string{"default1", "default2"}
		assert.Equal(t, defaultSlice, c.GetStringSlice("nonexistent", defaultSlice))
		assert.Equal(t, []string{}, c.GetStringSlice("nonexistent"))

		// Test conversion from other types - only handles slices and strings with commas
		assert.Equal(t, []string{"123"}, c.GetStringSlice("string_int")) // This should convert string to slice
	})

	t.Run("GetNestedMap", func(t *testing.T) {
		// Test getting nested map
		nestedMap := c.GetNestedMap("nested.level1")
		require.NotNil(t, nestedMap)
		assert.Equal(t, "nested_string", nestedMap["string_val"])
		assert.Equal(t, 100, nestedMap["int_val"])
		assert.Equal(t, true, nestedMap["bool_val"])

		// Test getting top-level map
		topMap := c.GetNestedMap("nested")
		require.NotNil(t, topMap)
		level1Map, ok := topMap["level1"].(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, "nested_string", level1Map["string_val"])

		// Test non-existent path
		assert.Nil(t, c.GetNestedMap("nonexistent.path"))

		// Test path that doesn't point to a map
		assert.Nil(t, c.GetNestedMap("string_value"))

		// Ensure returned map is a copy (thread safety)
		nestedMap["new_key"] = "new_value"
		freshMap := c.GetNestedMap("nested.level1")
		assert.NotContains(t, freshMap, "new_key")
	})
}

// TestGetters_EdgeCases tests edge cases and error conditions
func TestGetters_EdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with mixed nested and flat keys
	testData := map[string]any{
		"flat_key": "flat_value",
		"nested": map[string]any{
			"key": "nested_value",
		},
		"nested.flat": "this_should_be_flat",
	}

	c.LoadFromMap(testData)

	t.Run("FlatVsNestedKeyPriority", func(t *testing.T) {
		// Flat key should take priority over nested interpretation
		assert.Equal(t, "this_should_be_flat", c.GetString("nested.flat"))

		// But nested access should still work for actual nested keys
		assert.Equal(t, "nested_value", c.GetString("nested.key"))
	})

	t.Run("DeepNesting", func(t *testing.T) {
		deepData := map[string]any{
			"level1": map[string]any{
				"level2": map[string]any{
					"level3": map[string]any{
						"level4": map[string]any{
							"value": "deep_value",
						},
					},
				},
			},
		}

		c.LoadFromMap(deepData)
		assert.Equal(t, "deep_value", c.GetString("level1.level2.level3.level4.value"))
	})

	t.Run("InvalidNestedPath", func(t *testing.T) {
		// Try to access nested path where intermediate value is not a map
		c.Set("not_map", "string_value")
		assert.Equal(t, "", c.GetString("not_map.something"))
		assert.Equal(t, 0, c.GetInt("not_map.something"))
		assert.False(t, c.GetBool("not_map.something"))
	})

	t.Run("EmptyKey", func(t *testing.T) {
		assert.Equal(t, "", c.GetString(""))
		assert.Equal(t, 0, c.GetInt(""))
		assert.False(t, c.GetBool(""))
	})

	t.Run("NilValues", func(t *testing.T) {
		c.Set("nil_value", nil)
		assert.Equal(t, "<nil>", c.GetString("nil_value")) // fmt.Sprintf formats nil as <nil>
		assert.Equal(t, 0, c.GetInt("nil_value"))
		assert.False(t, c.GetBool("nil_value"))
	})

	// Test nil config edge case for all getters
	t.Run("NilConfig", func(t *testing.T) {
		var nilConfig *Config
		assert.NotPanics(t, func() {
			assert.Equal(t, "", nilConfig.GetString("any"))
			assert.Equal(t, 0, nilConfig.GetInt("any"))
			assert.Equal(t, 0.0, nilConfig.GetFloat64("any"))
			assert.False(t, nilConfig.GetBool("any"))
			assert.Equal(t, time.Duration(0), nilConfig.GetDuration("any"))
			assert.Nil(t, nilConfig.GetStringSlice("any"))
			assert.Empty(t, nilConfig.GetNestedMap("any"))
			assert.Empty(t, nilConfig.GetNestedKeys("any"))
			assert.Nil(t, nilConfig.GetAll())
		})
	})
}

// TestGetters_GetNestedKeys_EdgeCases tests edge cases for GetNestedKeys
func TestGetters_GetNestedKeys_EdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with empty config
	keys := c.GetNestedKeys("nonexistent")
	assert.Empty(t, keys)

	// Test with nil config
	var nilConfig *Config
	keys = nilConfig.GetNestedKeys("any")
	assert.Empty(t, keys)

	// Test with complex nested structure
	c.LoadFromMap(map[string]any{
		"level1": map[string]any{
			"level2a": map[string]any{
				"level3": "value",
			},
			"level2b": "simple_value",
		},
		"simple": "value",
	})

	keys = c.GetNestedKeys("level1")
	assert.Contains(t, keys, "level1.level2a")
	assert.Contains(t, keys, "level1.level2b")
	assert.Len(t, keys, 2) // Both level2a and level2b are immediate children
}

// TestGetters_GetAll_EdgeCases tests edge cases for GetAll
func TestGetters_GetAll_EdgeCases(t *testing.T) {
	// Test with nil config
	var c *Config
	all := c.GetAll()
	assert.Nil(t, all)

	// Test with empty config
	c, err := New()
	require.NoError(t, err)
	all = c.GetAll()
	assert.Empty(t, all)
}

// TestGetters_TypeConversion_Errors tests type conversion error cases
func TestGetters_TypeConversion_Errors(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	c.LoadFromMap(map[string]any{
		"invalid_duration": "not_a_duration",
		"complex_type":     complex(1, 2),
		"nil_value":        nil,
	})

	// Test invalid duration - should return default
	duration := c.GetDuration("invalid_duration", 5*time.Second)
	assert.Equal(t, 5*time.Second, duration)

	// Test with complex type
	str := c.GetString("complex_type")
	assert.Contains(t, str, "+")

	// Test with nil value
	str = c.GetString("nil_value")
	assert.Equal(t, "<nil>", str) // GetString with nil returns "<nil>" due to fmt.Sprintf
}

// Benchmark Tests for getters.go functions

func BenchmarkConfig_GetString(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"string_key": "test_value",
		"nested": map[string]any{
			"string_key": "nested_value",
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetString("string_key")
		c.GetString("nested.string_key")
		c.GetString("nonexistent", "default")
	}
}

func BenchmarkConfig_GetInt(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"int_key":        42,
		"string_int_key": "123",
		"nested": map[string]any{
			"int_key": 84,
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetInt("int_key")
		c.GetInt("string_int_key")
		c.GetInt("nested.int_key")
		c.GetInt("nonexistent", 0)
	}
}

func BenchmarkConfig_GetFloat64(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"float_key":        3.14159,
		"string_float_key": "2.71828",
		"int_key":          42,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetFloat64("float_key")
		c.GetFloat64("string_float_key")
		c.GetFloat64("int_key")
		c.GetFloat64("nonexistent", 0.0)
	}
}

func BenchmarkConfig_GetBool(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"bool_key":        true,
		"string_bool_key": "false",
		"int_bool_key":    1,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetBool("bool_key")
		c.GetBool("string_bool_key")
		c.GetBool("int_bool_key")
		c.GetBool("nonexistent", false)
	}
}

func BenchmarkConfig_GetDuration(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"duration_key":        "5m30s",
		"string_duration_key": "1h",
		"int_duration_key":    300,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetDuration("duration_key")
		c.GetDuration("string_duration_key")
		c.GetDuration("int_duration_key")
		c.GetDuration("nonexistent", time.Second)
	}
}

func BenchmarkConfig_GetStringSlice(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"slice_key":       []string{"a", "b", "c"},
		"string_key":      "x,y,z",
		"interface_slice": []interface{}{"d", "e", "f"},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetStringSlice("slice_key")
		c.GetStringSlice("string_key")
		c.GetStringSlice("interface_slice")
		c.GetStringSlice("nonexistent", []string{"default"})
	}
}

func BenchmarkConfig_GetNestedMap(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"nested": map[string]any{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetNestedMap("nested")
		c.GetNestedMap("nonexistent")
	}
}

func BenchmarkConfig_GetNestedKeys(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"level1": map[string]any{
			"level2a": map[string]any{
				"level3a": "value1",
				"level3b": "value2",
			},
			"level2b": "value3",
			"level2c": map[string]any{
				"level3c": "value4",
			},
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetNestedKeys("level1")
		c.GetNestedKeys("nonexistent")
	}
}

func BenchmarkConfig_GetAll(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	// Add multiple keys
	data := make(map[string]any)
	for i := 0; i < 100; i++ {
		data[fmt.Sprintf("key%d", i)] = fmt.Sprintf("value%d", i)
	}
	c.LoadFromMap(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.GetAll()
	}
}
