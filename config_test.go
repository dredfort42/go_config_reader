/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | config_test.go
	::  ::          ::  ::    Created  | 2025-08-07
		  ::::  ::::          Modified | 2025-08-19

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
	assert.Equal(t, time.Duration(30.5*float64(time.Second)), c.GetDuration("duration_float"))

	c.Set("duration_int64", int64(45))
	assert.Equal(t, 45*time.Second, c.GetDuration("duration_int64"))

	c.Set("duration_invalid_string", "not_a_duration")
	assert.Equal(t, time.Duration(0), c.GetDuration("duration_invalid_string"))

	// Test GetStringSlice with []any
	c.Set("slice_any", []any{"a", 1, true})
	result := c.GetStringSlice("slice_any")
	assert.Equal(t, []string{"a", "1", "true"}, result)
}

func TestConfig_ParseINI_Sections(t *testing.T) {
	iniContent := `
# Global settings
global_key=global_value
debug=true

[server]
host=localhost
port=8080
ssl_enabled=true

[database]
host=db.example.com
port=5432
name=testdb
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

	// Test global keys
	assert.Equal(t, "global_value", c.GetString("global_key"))
	assert.True(t, c.GetBool("debug"))

	// Test nested section access
	assert.Equal(t, "localhost", c.GetString("server.host"))
	assert.Equal(t, 8080, c.GetInt("server.port"))
	assert.True(t, c.GetBool("server.ssl_enabled"))

	assert.Equal(t, "db.example.com", c.GetString("database.host"))
	assert.Equal(t, 5432, c.GetInt("database.port"))
	assert.Equal(t, "testdb", c.GetString("database.name"))
}

func TestConfig_ParseINI_QuotedValues(t *testing.T) {
	iniContent := `
[test]
double_quoted="value with spaces"
single_quoted='another value'
quoted_with_quotes="value with \"inner quotes\""
escape_sequences="line1\nline2\ttab"
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

	assert.Equal(t, "value with spaces", c.GetString("test.double_quoted"))
	assert.Equal(t, "another value", c.GetString("test.single_quoted"))
	assert.Equal(t, "value with \"inner quotes\"", c.GetString("test.quoted_with_quotes"))
	assert.Equal(t, "line1\nline2\ttab", c.GetString("test.escape_sequences"))
}

func TestConfig_ParseINI_BooleanValues(t *testing.T) {
	iniContent := `
[booleans]
true1=true
true2=yes
true3=on
true4=1
false1=false
false2=no
false3=off
false4=0
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

	// Test true values
	assert.True(t, c.GetBool("booleans.true1"))
	assert.True(t, c.GetBool("booleans.true2"))
	assert.True(t, c.GetBool("booleans.true3"))
	assert.True(t, c.GetBool("booleans.true4"))

	// Test false values
	assert.False(t, c.GetBool("booleans.false1"))
	assert.False(t, c.GetBool("booleans.false2"))
	assert.False(t, c.GetBool("booleans.false3"))
	assert.False(t, c.GetBool("booleans.false4"))
}

func TestConfig_ParseINI_NumericValues(t *testing.T) {
	iniContent := `
[numbers]
integer=42
negative=-273
float=3.14159
scientific=1.23e-4
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

	assert.Equal(t, 42, c.GetInt("numbers.integer"))
	assert.Equal(t, -273, c.GetInt("numbers.negative"))
	assert.Equal(t, 3.14159, c.GetFloat64("numbers.float"))
	assert.Equal(t, 1.23e-4, c.GetFloat64("numbers.scientific"))
}

func TestConfig_ParseINI_Lists(t *testing.T) {
	iniContent := `
[lists]
simple_list=item1,item2,item3
spaced_list=item 1, item 2, item 3
mixed_list=string,123,true,false
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

	simpleList := c.GetStringSlice("lists.simple_list")
	assert.Equal(t, []string{"item1", "item2", "item3"}, simpleList)

	spacedList := c.GetStringSlice("lists.spaced_list")
	assert.Equal(t, []string{"item 1", "item 2", "item 3"}, spacedList)

	mixedList := c.GetStringSlice("lists.mixed_list")
	assert.Equal(t, []string{"string", "123", "true", "false"}, mixedList)
}

func TestConfig_ParseINI_MultilineValues(t *testing.T) {
	iniContent := `
[multiline]
continued_value=line1 \
    line2 \
    line3
normal_value=single_line
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

	assert.Equal(t, "line1 line2 line3", c.GetString("multiline.continued_value"))
	assert.Equal(t, "single_line", c.GetString("multiline.normal_value"))
}

func TestConfig_ParseINI_Comments(t *testing.T) {
	iniContent := `
# This is a comment
; This is also a comment
key1=value1 # inline comment
key2=value2 ; inline comment
; key3=commented_out

[section]
key4=value4
# key5=commented_out
key6="quoted value # not a comment"
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

	assert.Equal(t, "value1", c.GetString("key1"))
	assert.Equal(t, "value2", c.GetString("key2"))
	assert.False(t, c.Has("key3"))

	assert.Equal(t, "value4", c.GetString("section.key4"))
	assert.False(t, c.Has("section.key5"))
	assert.Equal(t, "quoted value # not a comment", c.GetString("section.key6"))
}

func TestConfig_ParseINI_InvalidSections(t *testing.T) {
	iniContent := `
[valid_section]
key1=value1

[]
key2=value2

[invalid[section]]
key3=value3

[section_with=equals]
key4=value4

[valid_section2]
key5=value5
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

	// Valid sections should work
	assert.Equal(t, "value1", c.GetString("valid_section.key1"))
	assert.Equal(t, "value5", c.GetString("valid_section2.key5"))

	// Invalid sections should be ignored
	assert.False(t, c.Has("key2")) // Empty section name
	assert.False(t, c.Has("key3")) // Invalid characters
	assert.False(t, c.Has("key4")) // Invalid characters
}

func TestConfig_ParseINI_EmptyValues(t *testing.T) {
	iniContent := `
[test]
empty_value=
whitespace_only=   
quoted_empty=""
quoted_whitespace="   "
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

	assert.Equal(t, "", c.GetString("test.empty_value"))
	assert.Equal(t, "", c.GetString("test.whitespace_only"))
	assert.Equal(t, "", c.GetString("test.quoted_empty"))
	assert.Equal(t, "   ", c.GetString("test.quoted_whitespace"))
}

func TestConfig_ParseINI_NestedMapAccess(t *testing.T) {
	iniContent := `
[database]
host=localhost
port=5432

[cache]
enabled=true
ttl=3600
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

	// Test nested map access
	dbMap := c.GetNestedMap("database")
	assert.NotNil(t, dbMap)
	assert.Equal(t, "localhost", dbMap["host"])
	assert.Equal(t, 5432, dbMap["port"])

	cacheMap := c.GetNestedMap("cache")
	assert.NotNil(t, cacheMap)
	assert.Equal(t, true, cacheMap["enabled"])
	assert.Equal(t, 3600, cacheMap["ttl"])

	// Test nested keys
	dbKeys := c.GetNestedKeys("database")
	assert.Contains(t, dbKeys, "database.host")
	assert.Contains(t, dbKeys, "database.port")
}

// TestConfig_Has_EdgeCases tests edge cases for Has method
func TestConfig_Has_EdgeCases(t *testing.T) {
	// Test with nil config
	var c *Config
	assert.False(t, c.Has("any_key"))

	// Test with empty config
	c, err := New()
	require.NoError(t, err)
	assert.False(t, c.Has(""))
	assert.False(t, c.Has("nonexistent"))

	// Test with nested keys
	c.LoadFromMap(map[string]any{
		"flat_key": "value",
		"nested": map[string]any{
			"key": "value",
		},
	})
	assert.True(t, c.Has("flat_key"))
	assert.True(t, c.Has("nested.key"))
	assert.False(t, c.Has("nested.nonexistent"))
}

// TestConfig_Size_EdgeCases tests edge cases for Size method
func TestConfig_Size_EdgeCases(t *testing.T) {
	// Test with nil config
	var c *Config
	assert.Equal(t, 0, c.Size())

	// Test with empty config
	c, err := New()
	require.NoError(t, err)
	assert.Equal(t, 0, c.Size())

	// Test after adding and removing data
	c.Set("key", "value")
	assert.Equal(t, 1, c.Size())
	c.Clear()
	assert.Equal(t, 0, c.Size())
}

// TestConfig_IsEmpty_EdgeCases tests edge cases for IsEmpty method
func TestConfig_IsEmpty_EdgeCases(t *testing.T) {
	// Test with nil config
	var c *Config
	assert.True(t, c.IsEmpty())

	// Test with empty config
	c, err := New()
	require.NoError(t, err)
	assert.True(t, c.IsEmpty())

	// Test after adding data
	c.Set("key", "value")
	assert.False(t, c.IsEmpty())

	// Test after clearing
	c.Clear()
	assert.True(t, c.IsEmpty())
}

// TestConfig_String_EdgeCases tests edge cases for String method
func TestConfig_String_EdgeCases(t *testing.T) {
	// Test with nil config
	var c *Config
	assert.Equal(t, "Config is nil", c.String())

	// Test with empty config
	c, err := New()
	require.NoError(t, err)
	assert.Equal(t, "", c.String())
}

// TestConfig_LoadFromFile_ErrorCases tests additional error cases
func TestConfig_LoadFromFile_ErrorCases(t *testing.T) {
	c, err := New()
	require.NoError(t, err)

	// Test with nonexistent file
	err = c.LoadFromFile("nonexistent.json", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nonexistent.json")

	// Test with invalid JSON
	tempFile, err := os.CreateTemp("", "invalid_*.json")
	require.NoError(t, err)
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write([]byte("{invalid json"))
	require.NoError(t, err)
	tempFile.Close()

	err = c.LoadFromFile(tempFile.Name(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON config")
}

// Benchmark Tests for config.go functions

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := New()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConfig_LoadFromFile_JSON(b *testing.B) {
	// Create a test JSON file
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
	if err != nil {
		b.Fatal(err)
	}

	tempFile, err := os.CreateTemp("", "bench_config_*.json")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(jsonData)
	if err != nil {
		b.Fatal(err)
	}
	tempFile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, err := New()
		if err != nil {
			b.Fatal(err)
		}
		err = c.LoadFromFile(tempFile.Name(), nil)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConfig_LoadFromMap(b *testing.B) {
	data := map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
		"nested": map[string]any{
			"subkey": "subvalue",
		},
	}

	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.LoadFromMap(data)
	}
}

func BenchmarkConfig_Has(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"existing_key": "value",
		"nested": map[string]any{
			"key": "value",
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Has("existing_key")
		c.Has("nonexistent_key")
		c.Has("nested.key")
	}
}

func BenchmarkConfig_Keys(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	// Add multiple keys
	for i := 0; i < 100; i++ {
		c.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Keys()
	}
}

func BenchmarkConfig_Size(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Size()
	}
}

func BenchmarkConfig_IsEmpty(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.IsEmpty()
	}
}

func BenchmarkConfig_Clear(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	data := map[string]any{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.LoadFromMap(data)
		c.Clear()
	}
}

func BenchmarkConfig_String(b *testing.B) {
	c, err := New()
	if err != nil {
		b.Fatal(err)
	}

	c.LoadFromMap(map[string]any{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.String()
	}
}
