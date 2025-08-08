/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | config_test.go
	::  ::          ::  ::    Created  | 2025-08-07
		  ::::  ::::          Modified | 2025-08-07

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c, err := New()
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, 0, len(c.Keys()))
}

func TestConfig_LoadFromFile_JSON(t *testing.T) {
	// Create a temporary JSON config file
	configData := map[string]interface{}{
		"server_port": 8080,
		"debug_mode":  true,
		"timeout":     30.5,
		"app_name":    "test-app",
		"features":    []string{"auth", "api", "web"},
		"database": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	// Load the configuration
	c, err := New()
	require.NoError(t, err)

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.NoError(t, err)

	// Test various getters
	assert.Equal(t, 8080, c.GetInt("server_port"))
	assert.True(t, c.GetBool("debug_mode"))
	assert.Equal(t, 30.5, c.GetFloat64("timeout"))
	assert.Equal(t, "test-app", c.GetString("app_name"))
}

func TestConfig_LoadFromFile_YAML(t *testing.T) {
	yamlContent := `
server_port: 9090
debug_mode: false
timeout: 45
app_name: "yaml-app"
features:
  - auth
  - api
database:
  host: "db.example.com"
  port: 3306
`

	tempFile, err := os.CreateTemp("", "config_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(yamlContent))
	require.NoError(t, err)
	tempFile.Close()

	c, err := New()
	require.NoError(t, err)

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 9090, c.GetInt("server_port"))
	assert.False(t, c.GetBool("debug_mode"))
	assert.Equal(t, 45.0, c.GetFloat64("timeout"))
	assert.Equal(t, "yaml-app", c.GetString("app_name"))
}

func TestConfig_LoadFromFile_INI(t *testing.T) {
	iniContent := `
# Server configuration
server_port=7070
debug_mode=true
timeout=25
app_name=ini-app

# Database settings
db_host=localhost
db_port=5432
`

	tempFile, err := os.CreateTemp("", "config_*.ini")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(iniContent))
	require.NoError(t, err)
	tempFile.Close()

	c, err := New()
	require.NoError(t, err)

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.NoError(t, err)

	assert.Equal(t, 7070, c.GetInt("server_port"))
	assert.True(t, c.GetBool("debug_mode"))
	assert.Equal(t, 25.0, c.GetFloat64("timeout"))
	assert.Equal(t, "ini-app", c.GetString("app_name"))
	assert.Equal(t, "localhost", c.GetString("db_host"))
	assert.Equal(t, 5432, c.GetInt("db_port"))
}

func TestConfig_EnvironmentOverride(t *testing.T) {
	// Set environment variable
	originalValue := os.Getenv("TEST_SERVER_PORT")
	os.Setenv("TEST_SERVER_PORT", "9999")
	defer func() {
		if originalValue == "" {
			os.Unsetenv("TEST_SERVER_PORT")
		} else {
			os.Setenv("TEST_SERVER_PORT", originalValue)
		}
	}()

	// Create config file
	configData := map[string]interface{}{
		"TEST_SERVER_PORT": 8080,
		"debug_mode":       true,
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	c, err := New()
	require.NoError(t, err)

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.NoError(t, err)

	// Environment variable should override file value
	assert.Equal(t, 9999, c.GetInt("TEST_SERVER_PORT"))
	assert.True(t, c.GetBool("debug_mode"))
}

func TestConfig_DefaultValues(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with defaults
	assert.Equal(t, "default", c.GetString("non_existent_key", "default"))
	assert.Equal(t, 42, c.GetInt("non_existent_key", 42))
	assert.True(t, c.GetBool("non_existent_key", true))
	assert.Equal(t, 3.14, c.GetFloat64("non_existent_key", 3.14))
	assert.Equal(t, 5*time.Second, c.GetDuration("non_existent_key", 5*time.Second))

	// Test without defaults
	assert.Equal(t, "", c.GetString("non_existent_key"))
	assert.Equal(t, 0, c.GetInt("non_existent_key"))
	assert.False(t, c.GetBool("non_existent_key"))
	assert.Equal(t, 0.0, c.GetFloat64("non_existent_key"))
	assert.Equal(t, time.Duration(0), c.GetDuration("non_existent_key"))
}

func TestConfig_LoadWithOptions(t *testing.T) {
	configData := map[string]interface{}{
		"server_port": 8080,
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	c, err := New()
	require.NoError(t, err)

	opts := &LoadOptions{
		DefaultValues: map[string]interface{}{
			"debug_mode": true,
			"timeout":    30,
		},
		RequiredKeys: []string{"server_port"},
	}

	err = c.LoadFromFile(tempFile.Name(), opts)
	assert.NoError(t, err)

	assert.Equal(t, 8080, c.GetInt("server_port"))
	assert.True(t, c.GetBool("debug_mode"))
	assert.Equal(t, 30, c.GetInt("timeout"))
}

func TestConfig_RequiredKeys(t *testing.T) {
	configData := map[string]interface{}{
		"server_port": 8080,
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	c, err := New()
	require.NoError(t, err)

	opts := &LoadOptions{
		RequiredKeys: []string{"server_port", "missing_key"},
	}

	err = c.LoadFromFile(tempFile.Name(), opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing_key")
}

func TestConfig_Validation(t *testing.T) {
	configData := map[string]interface{}{
		"server_port": 80, // Port < 1024 should fail validation
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	c, err := New()
	require.NoError(t, err)

	opts := &LoadOptions{
		ValidationFunc: func(data map[string]interface{}) error {
			if port, ok := data["server_port"]; ok {
				var portInt int
				switch v := port.(type) {
				case int:
					portInt = v
				case float64:
					portInt = int(v)
				default:
					return fmt.Errorf("server_port must be a number")
				}
				if portInt < 1024 {
					return fmt.Errorf("server port must be >= 1024")
				}
			}

			return nil
		},
	}

	// This should fail validation since port is 80 (< 1024)
	err = c.LoadFromFile(tempFile.Name(), opts)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestConfig_ThreadSafety(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	c.Set("test_key", "test_value")

	// Test concurrent access
	done := make(chan bool, 2)

	go func() {
		for i := 0; i < 100; i++ {
			c.Set("concurrent_key", i)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			_ = c.GetInt("concurrent_key", 0)
		}
		done <- true
	}()

	<-done
	<-done

	assert.True(t, c.Has("test_key"))
	assert.Equal(t, "test_value", c.GetString("test_key"))
}

func TestConfig_GetStringSlice(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with actual slice
	c.Set("features", []string{"auth", "api", "web"})
	features := c.GetStringSlice("features")
	assert.Equal(t, []string{"auth", "api", "web"}, features)

	// Test with comma-separated string
	c.Set("tags", "tag1,tag2,tag3")
	tags := c.GetStringSlice("tags")
	assert.Equal(t, []string{"tag1", "tag2", "tag3"}, tags)

	// Test with default
	defaultSlice := []string{"default1", "default2"}
	result := c.GetStringSlice("non_existent", defaultSlice)
	assert.Equal(t, defaultSlice, result)

	// Test empty
	empty := c.GetStringSlice("non_existent")
	assert.Equal(t, []string{}, empty)
}

func TestConfig_GetDuration(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test duration string
	c.Set("timeout", "30s")
	assert.Equal(t, 30*time.Second, c.GetDuration("timeout"))

	// Test integer (seconds)
	c.Set("delay", 60)
	assert.Equal(t, 60*time.Second, c.GetDuration("delay"))

	// Test complex duration
	c.Set("long_timeout", "1h30m45s")
	expected := 1*time.Hour + 30*time.Minute + 45*time.Second
	assert.Equal(t, expected, c.GetDuration("long_timeout"))

	// Test default
	defaultDuration := 5 * time.Minute
	result := c.GetDuration("non_existent", defaultDuration)
	assert.Equal(t, defaultDuration, result)
}

func TestConfig_Keys(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	assert.Equal(t, 0, len(c.Keys()))

	c.Set("key1", "value1")
	c.Set("key2", "value2")

	keys := c.Keys()
	assert.Equal(t, 2, len(keys))
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestConfig_GetAll(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	c.Set("key1", "value1")
	c.Set("key2", 42)

	all := c.GetAll()
	assert.Equal(t, 2, len(all))
	assert.Equal(t, "value1", all["key1"])
	assert.Equal(t, 42, all["key2"])

	// Verify it's a copy
	all["key3"] = "value3"
	assert.False(t, c.Has("key3"))
}

func TestGlobalAPI(t *testing.T) {
	// Test the global API functions
	Set("global_key", "global_value")
	assert.True(t, Has("global_key"))
	assert.Equal(t, "global_value", GetString("global_key"))

	Set("global_int", 123)
	assert.Equal(t, 123, GetInt("global_int"))

	Set("global_bool", true)
	assert.True(t, GetBool("global_bool"))

	keys := Keys()
	assert.Contains(t, keys, "global_key")
	assert.Contains(t, keys, "global_int")
	assert.Contains(t, keys, "global_bool")
}

func TestMustLoad(t *testing.T) {
	configData := map[string]interface{}{
		"server_port": 8080,
		"debug_mode":  true,
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	// Should not panic with valid config
	c := MustLoad(tempFile.Name(), nil)
	assert.NotNil(t, c)
	assert.Equal(t, 8080, c.GetInt("server_port"))

	// Should panic with invalid file
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected MustLoad to panic with invalid file")
		}
	}()
	MustLoad("/non/existent/file.json", nil)
}

func TestLoadWithDefaults(t *testing.T) {
	configData := map[string]interface{}{
		"server_port": 8080,
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	defaults := map[string]interface{}{
		"debug_mode":  true,
		"timeout":     30,
		"server_port": 9090, // Should be overridden by file
	}

	c := LoadWithDefaults(tempFile.Name(), defaults)
	assert.Equal(t, 8080, c.GetInt("server_port")) // From file
	assert.True(t, c.GetBool("debug_mode"))        // From defaults
	assert.Equal(t, 30, c.GetInt("timeout"))       // From defaults
}

func BenchmarkConfig_GetString(b *testing.B) {
	c, _ := New()
	c.Set("test_key", "test_value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.GetString("test_key")
	}
}

func BenchmarkConfig_GetInt(b *testing.B) {
	c, _ := New()
	c.Set("test_key", 42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = c.GetInt("test_key")
	}
}

// Tests for compatibility.go functions with 0% coverage
func TestGlobalGetFloat64(t *testing.T) {
	Set("float_key", 3.14)
	assert.Equal(t, 3.14, GetFloat64("float_key"))

	Set("float_string", "2.71")
	assert.Equal(t, 2.71, GetFloat64("float_string"))

	assert.Equal(t, 1.5, GetFloat64("non_existent", 1.5))
	assert.Equal(t, 0.0, GetFloat64("non_existent"))
}

func TestGlobalGetDuration(t *testing.T) {
	Set("duration_key", "30s")
	assert.Equal(t, 30*time.Second, GetDuration("duration_key"))

	Set("duration_int", 60)
	assert.Equal(t, 60*time.Second, GetDuration("duration_int"))

	defaultDur := 5 * time.Minute
	assert.Equal(t, defaultDur, GetDuration("non_existent", defaultDur))
	assert.Equal(t, time.Duration(0), GetDuration("non_existent"))
}

func TestGlobalGetStringSlice(t *testing.T) {
	Set("slice_key", []string{"a", "b", "c"})
	assert.Equal(t, []string{"a", "b", "c"}, GetStringSlice("slice_key"))

	Set("csv_key", "x,y,z")
	assert.Equal(t, []string{"x", "y", "z"}, GetStringSlice("csv_key"))

	defaultSlice := []string{"default1", "default2"}
	assert.Equal(t, defaultSlice, GetStringSlice("non_existent", defaultSlice))
	assert.Equal(t, []string{}, GetStringSlice("non_existent"))
}

func TestGlobalLoad(t *testing.T) {
	configData := map[string]interface{}{
		"test_key": "test_value",
		"port":     8080,
	}

	jsonData, err := json.Marshal(configData)
	require.NoError(t, err)

	tempFile, err := os.CreateTemp("", "config_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	require.NoError(t, err)
	tempFile.Close()

	err = Load(tempFile.Name(), nil)
	assert.NoError(t, err)

	assert.Equal(t, "test_value", GetString("test_key"))
	assert.Equal(t, 8080, GetInt("port"))

	// Test Load with options
	opts := &LoadOptions{
		DefaultValues: map[string]interface{}{
			"default_key": "default_value",
		},
	}
	err = Load(tempFile.Name(), opts)
	assert.NoError(t, err)
	assert.Equal(t, "default_value", GetString("default_key"))
}

func TestGlobalString(t *testing.T) {
	// Clear any existing global config and set some test values
	Set("key1", "value1")
	Set("key2", 42)

	str := String()
	assert.NotEmpty(t, str)
	assert.Contains(t, str, "key1: value1")
	assert.Contains(t, str, "key2: 42")
}

// Tests for config.go functions with 0% coverage
func TestConfig_LoadFromMap(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	data := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	c.LoadFromMap(data)

	assert.Equal(t, "value1", c.GetString("key1"))
	assert.Equal(t, 42, c.GetInt("key2"))
	assert.True(t, c.GetBool("key3"))
	assert.Equal(t, 3, c.Size())
}

func TestConfig_Clear(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	c.Set("key1", "value1")
	c.Set("key2", "value2")
	assert.Equal(t, 2, c.Size())
	assert.False(t, c.IsEmpty())

	c.Clear()
	assert.Equal(t, 0, c.Size())
	assert.True(t, c.IsEmpty())
	assert.False(t, c.Has("key1"))
	assert.False(t, c.Has("key2"))
}

func TestConfig_Size(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	assert.Equal(t, 0, c.Size())

	c.Set("key1", "value1")
	assert.Equal(t, 1, c.Size())

	c.Set("key2", "value2")
	assert.Equal(t, 2, c.Size())

	c.Clear()
	assert.Equal(t, 0, c.Size())
}

func TestConfig_IsEmpty(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	assert.True(t, c.IsEmpty())

	c.Set("key1", "value1")
	assert.False(t, c.IsEmpty())

	c.Clear()
	assert.True(t, c.IsEmpty())
}

func TestConfig_String(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	c.Set("key1", "value1")
	c.Set("key2", 42)

	str := c.String()
	assert.NotEmpty(t, str)
	assert.Contains(t, str, "key1: value1")
	assert.Contains(t, str, "key2: 42")

	// Test empty config
	c.Clear()
	str = c.String()
	assert.Equal(t, "", str)
}

// Tests for low coverage functions
func TestNewWithOptions(t *testing.T) {
	// Test New with invalid option
	invalidOption := func(c *Config) error {
		return fmt.Errorf("test error")
	}

	_, err := New(invalidOption)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to apply option")

	// Test New with valid option
	validOption := func(c *Config) error {
		c.Set("option_key", "option_value")

		return nil
	}

	c, err := New(validOption)
	assert.NoError(t, err)
	assert.Equal(t, "option_value", c.GetString("option_key"))
}

func TestInitGlobalConfigPanic(t *testing.T) {
	// This is hard to test since it would require causing New() to fail
	// and we can't easily reset the sync.Once. The current implementation
	// should not panic under normal circumstances, but we can at least
	// ensure initGlobalConfig doesn't panic when called multiple times
	initGlobalConfig()
	initGlobalConfig() // Should not panic on second call
}

// Additional parseINI edge cases to improve coverage
func TestParseINIEdgeCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test INI parsing with various edge cases
	iniContent := `
# Comment line
; Another comment
key_with_spaces = value with spaces
quoted_value = "quoted string"
single_quoted = 'single quoted'
key_with_inline_comment = value # inline comment
key_with_semicolon_comment = value ; semicolon comment
empty_value = 
key_with_equals_in_value = value=with=equals
key without equals sign - should be ignored
= value_without_key
`

	tempFile, err := os.CreateTemp("", "config_*.ini")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(iniContent))
	require.NoError(t, err)
	tempFile.Close()

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.NoError(t, err)

	assert.Equal(t, "value with spaces", c.GetString("key_with_spaces"))
	assert.Equal(t, "quoted string", c.GetString("quoted_value"))
	assert.Equal(t, "single quoted", c.GetString("single_quoted"))
	assert.Equal(t, "value", c.GetString("key_with_inline_comment"))
	assert.Equal(t, "value", c.GetString("key_with_semicolon_comment"))
	assert.Equal(t, "", c.GetString("empty_value"))
	assert.Equal(t, "value=with=equals", c.GetString("key_with_equals_in_value"))
	assert.False(t, c.Has("key without equals sign"))
	assert.False(t, c.Has(""))
}

// Test error cases for better coverage
func TestConfig_ErrorCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test loading non-existent file
	err = c.LoadFromFile("/non/existent/file.json", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration file not found")

	// Test invalid JSON
	invalidJSON := `{"invalid": json}`
	tempFile, err := os.CreateTemp("", "invalid_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte(invalidJSON))
	require.NoError(t, err)
	tempFile.Close()

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON config")

	// Test invalid YAML
	invalidYAML := `invalid: yaml: content: [`
	tempYAMLFile, err := os.CreateTemp("", "invalid_*.yaml")
	require.NoError(t, err)
	defer os.Remove(tempYAMLFile.Name())

	_, err = tempYAMLFile.Write([]byte(invalidYAML))
	require.NoError(t, err)
	tempYAMLFile.Close()

	err = c.LoadFromFile(tempYAMLFile.Name(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse YAML config")
}

// Test type conversion edge cases for better coverage
func TestConfig_TypeConversions(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test GetFloat64 with different types
	c.Set("float_int", 42)
	assert.Equal(t, 42.0, c.GetFloat64("float_int"))

	c.Set("float_int64", int64(123))
	assert.Equal(t, 123.0, c.GetFloat64("float_int64"))

	c.Set("float_float32", float32(3.14))
	assert.InDelta(t, 3.14, c.GetFloat64("float_float32"), 0.01)

	c.Set("float_invalid_string", "not_a_number")
	assert.Equal(t, 0.0, c.GetFloat64("float_invalid_string"))

	// Test GetInt with int64
	c.Set("int_int64", int64(456))
	assert.Equal(t, 456, c.GetInt("int_int64"))

	c.Set("int_invalid_string", "not_a_number")
	assert.Equal(t, 0, c.GetInt("int_invalid_string"))

	// Test GetDuration with float64
	c.Set("duration_float", 30.5)
	assert.Equal(t, 30*time.Second, c.GetDuration("duration_float"))

	c.Set("duration_int64", int64(45))
	assert.Equal(t, 45*time.Second, c.GetDuration("duration_int64"))

	c.Set("duration_invalid_string", "not_a_duration")
	assert.Equal(t, time.Duration(0), c.GetDuration("duration_invalid_string"))

	// Test GetStringSlice with []any
	c.Set("slice_any", []any{"a", 1, true})
	result := c.GetStringSlice("slice_any")
	assert.Equal(t, []string{"a", "1", "true"}, result)
}

func BenchmarkConfig_LoadFromFile(b *testing.B) {
	configData := map[string]interface{}{
		"server_port": 8080,
		"debug_mode":  true,
		"timeout":     30,
	}

	jsonData, _ := json.Marshal(configData)
	tempFile, _ := os.CreateTemp("", "config_*.json")
	defer os.Remove(tempFile.Name())
	tempFile.Write(jsonData)
	tempFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, _ := New()
		_ = c.LoadFromFile(tempFile.Name(), nil)
	}
}
