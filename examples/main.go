/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | main.go
	::  ::          ::  ::    Created  | 2025-08-19
		  ::::  ::::          Modified | 2025-08-19

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	config "github.com/dredfort42/go_config_reader"
)

func main() {
	fmt.Println("🔧 Go Configuration Reader - Comprehensive Examples")
	fmt.Println("===================================================\n")

	// Example 1: Basic Usage
	fmt.Println("1. Basic Configuration Loading")
	demonstrateBasicUsage()

	// Example 2: Different File Formats
	fmt.Println("2. Different Configuration Formats")
	demonstrateFileFormats()

	// Example 3: Type Safety and Getters
	fmt.Println("3. Type Safety and Getter Methods")
	demonstrateTypeSafety()

	// Example 4: Nested Configuration
	fmt.Println("4. Nested Configuration with Dot Notation")
	demonstrateNestedConfiguration()

	// Example 5: Default Values
	fmt.Println("5. Default Values")
	demonstrateDefaultValues()

	// Example 6: Environment Variable Override
	fmt.Println("6. Environment Variable Override")
	demonstrateEnvironmentOverride()

	// Example 7: Advanced Loading Options
	fmt.Println("7. Advanced Loading with Options")
	demonstrateAdvancedOptions()

	// Example 8: Runtime Configuration Changes
	fmt.Println("8. Runtime Configuration Changes")
	demonstrateRuntimeChanges()

	// Example 9: Utility Methods
	fmt.Println("9. Utility Methods")
	demonstrateUtilityMethods()

	// Example 10: Error Handling
	fmt.Println("10. Error Handling")
	demonstrateErrorHandling()

	// Example 11: Thread Safety
	fmt.Println("11. Thread Safety")
	demonstrateThreadSafety()

	// Example 12: Complex Real-World Example
	fmt.Println("12. Real-World Application Configuration")
	demonstrateRealWorldExample()
}

// Example 1: Basic Usage
func demonstrateBasicUsage() {
	// Create a new config instance
	cfg, err := config.New()
	if err != nil {
		log.Printf("Error creating config: %v\n", err)
		return
	}

	// Load from map (simulating loaded config)
	cfg.LoadFromMap(map[string]any{
		"app_name":    "MyApplication",
		"server_port": 8080,
		"debug_mode":  true,
		"timeout":     "30s",
	})

	// Get values with automatic type conversion
	appName := cfg.GetString("app_name", "DefaultApp")
	port := cfg.GetInt("server_port", 3000)
	debug := cfg.GetBool("debug_mode", false)
	timeout := cfg.GetDuration("timeout", 10*time.Second)

	fmt.Printf("  Application: %s\n", appName)
	fmt.Printf("  Port: %d\n", port)
	fmt.Printf("  Debug: %v\n", debug)
	fmt.Printf("  Timeout: %v\n", timeout)
	fmt.Println()
}

// Example 2: Different File Formats
func demonstrateFileFormats() {
	// Create sample configuration files
	createSampleFiles()

	formats := []struct {
		name string
		file string
	}{
		{"JSON", "config.json"},
		{"YAML", "config.yaml"},
		{"INI", "config.ini"},
	}

	for _, format := range formats {
		fmt.Printf("  Loading %s configuration:\n", format.name)

		cfg, err := config.New()
		if err != nil {
			fmt.Printf("    Error creating config: %v\n", err)
			continue
		}

		err = cfg.LoadFromFile(format.file, nil)
		if err != nil {
			fmt.Printf("    Error loading %s: %v\n", format.file, err)
			continue
		}

		fmt.Printf("    Loaded %d keys successfully\n", cfg.Size())

		// Try to get a common value
		if cfg.Has("server_port") {
			fmt.Printf("    Server Port: %d\n", cfg.GetInt("server_port"))
		}
		if cfg.Has("app_name") {
			fmt.Printf("    App Name: %s\n", cfg.GetString("app_name"))
		}
	}
	fmt.Println()
}

// Example 3: Type Safety and Getters
func demonstrateTypeSafety() {
	cfg, _ := config.New()

	// Load various data types
	cfg.LoadFromMap(map[string]any{
		"string_value":    "hello world",
		"int_value":       42,
		"float_value":     3.14159,
		"bool_value":      true,
		"duration_string": "1h30m",
		"duration_int":    300, // seconds
		"string_slice":    []string{"apple", "banana", "cherry"},
		"mixed_slice":     []any{"red", "green", "blue"},
		"comma_separated": "item1,item2,item3",
	})

	fmt.Println("  Type-safe value retrieval:")
	fmt.Printf("    String: %s\n", cfg.GetString("string_value"))
	fmt.Printf("    Integer: %d\n", cfg.GetInt("int_value"))
	fmt.Printf("    Float: %.2f\n", cfg.GetFloat64("float_value"))
	fmt.Printf("    Boolean: %v\n", cfg.GetBool("bool_value"))
	fmt.Printf("    Duration (string): %v\n", cfg.GetDuration("duration_string"))
	fmt.Printf("    Duration (int): %v\n", cfg.GetDuration("duration_int"))
	fmt.Printf("    String Slice: %v\n", cfg.GetStringSlice("string_slice"))
	fmt.Printf("    Mixed Slice: %v\n", cfg.GetStringSlice("mixed_slice"))
	fmt.Printf("    Comma Separated: %v\n", cfg.GetStringSlice("comma_separated"))

	// Demonstrate type conversion
	fmt.Println("  Automatic type conversion:")
	cfg.Set("number_as_string", "123")
	cfg.Set("bool_as_string", "true")
	fmt.Printf("    String to Int: %d\n", cfg.GetInt("number_as_string"))
	fmt.Printf("    String to Bool: %v\n", cfg.GetBool("bool_as_string"))
	fmt.Println()
}

// Example 4: Nested Configuration
func demonstrateNestedConfiguration() {
	cfg, _ := config.New()

	// Load nested configuration
	cfg.LoadFromMap(map[string]any{
		"database": map[string]any{
			"host":    "localhost",
			"port":    5432,
			"name":    "myapp",
			"ssl":     true,
			"timeout": "30s",
			"pool": map[string]any{
				"min_connections": 5,
				"max_connections": 20,
			},
		},
		"redis": map[string]any{
			"host": "redis.example.com",
			"port": 6379,
		},
	})

	fmt.Println("  Nested configuration access with dot notation:")
	fmt.Printf("    Database Host: %s\n", cfg.GetString("database.host"))
	fmt.Printf("    Database Port: %d\n", cfg.GetInt("database.port"))
	fmt.Printf("    Database SSL: %v\n", cfg.GetBool("database.ssl"))
	fmt.Printf("    Database Timeout: %v\n", cfg.GetDuration("database.timeout"))
	fmt.Printf("    Pool Min Connections: %d\n", cfg.GetInt("database.pool.min_connections"))
	fmt.Printf("    Pool Max Connections: %d\n", cfg.GetInt("database.pool.max_connections"))
	fmt.Printf("    Redis Host: %s\n", cfg.GetString("redis.host"))

	// Get nested map
	dbConfig := cfg.GetNestedMap("database")
	if dbConfig != nil {
		fmt.Printf("    Database config map keys: %v\n", getMapKeys(dbConfig))
	}

	// Get nested keys
	dbKeys := cfg.GetNestedKeys("database")
	fmt.Printf("    Database nested keys: %v\n", dbKeys)

	// Set nested values at runtime
	cfg.Set("database.pool.idle_timeout", "5m")
	fmt.Printf("    Set nested value - Idle Timeout: %v\n", cfg.GetDuration("database.pool.idle_timeout"))
	fmt.Println()
}

// Example 5: Default Values
func demonstrateDefaultValues() {
	cfg, _ := config.New()

	// Load minimal configuration
	cfg.LoadFromMap(map[string]any{
		"app_name": "ProductionApp",
	})

	fmt.Println("  Using default values for missing keys:")
	fmt.Printf("    App Name: %s\n", cfg.GetString("app_name", "DefaultApp"))
	fmt.Printf("    Server Port (default): %d\n", cfg.GetInt("server_port", 8080))
	fmt.Printf("    Debug Mode (default): %v\n", cfg.GetBool("debug_mode", false))
	fmt.Printf("    Timeout (default): %v\n", cfg.GetDuration("timeout", 30*time.Second))
	fmt.Printf("    Features (default): %v\n", cfg.GetStringSlice("features", []string{"basic", "auth"}))
	fmt.Println()
}

// Example 6: Environment Variable Override
func demonstrateEnvironmentOverride() {
	// Set some environment variables
	os.Setenv("server_port", "9090")
	os.Setenv("debug_mode", "true")
	os.Setenv("app_name", "EnvOverrideApp")
	defer func() {
		os.Unsetenv("server_port")
		os.Unsetenv("debug_mode")
		os.Unsetenv("app_name")
	}()

	cfg, _ := config.New()

	// Load configuration that will be overridden
	cfg.LoadFromMap(map[string]any{
		"server_port": 8080,
		"debug_mode":  false,
		"app_name":    "FileApp",
	})

	// Since we loaded from map, manually apply environment override
	for key := range cfg.GetAll() {
		if envValue := os.Getenv(key); envValue != "" {
			cfg.Set(key, envValue)
		}
	}

	fmt.Println("  Environment variable override:")
	fmt.Printf("    Server Port (from ENV): %d\n", cfg.GetInt("server_port"))
	fmt.Printf("    Debug Mode (from ENV): %v\n", cfg.GetBool("debug_mode"))
	fmt.Printf("    App Name (from ENV): %s\n", cfg.GetString("app_name"))

	// Example with environment override disabled
	cfg2, _ := config.New()
	cfg2.LoadFromMap(map[string]any{
		"server_port": 8080,
		"debug_mode":  false,
		"app_name":    "FileApp",
	})

	fmt.Println("  Without environment override:")
	fmt.Printf("    Server Port (from config): %d\n", cfg2.GetInt("server_port"))
	fmt.Printf("    Debug Mode (from config): %v\n", cfg2.GetBool("debug_mode"))
	fmt.Printf("    App Name (from config): %s\n", cfg2.GetString("app_name"))
	fmt.Println()
}

// Example 7: Advanced Loading Options
func demonstrateAdvancedOptions() {
	cfg, _ := config.New()

	// Advanced loading options
	opts := &config.LoadOptions{
		DefaultValues: map[string]any{
			"server_port":     8080,
			"debug_mode":      false,
			"max_connections": 100,
			"timeout":         "30s",
		},
		RequiredKeys: []string{"app_name", "database_url"},
		ValidationFunc: func(data map[string]any) error {
			// Custom validation logic
			if port, ok := data["server_port"].(int); ok && port < 1024 {
				return fmt.Errorf("server port must be >= 1024, got %d", port)
			}
			if timeout, ok := data["timeout"].(string); ok {
				if _, err := time.ParseDuration(timeout); err != nil {
					return fmt.Errorf("invalid timeout format: %s", timeout)
				}
			}
			return nil
		},
	}

	// Create a valid configuration
	validConfig := map[string]any{
		"app_name":     "ValidatedApp",
		"database_url": "postgres://localhost/myapp",
		"server_port":  8080,
	}

	// Save to temporary file and load with options
	tempFile := createTempConfig(validConfig)
	defer os.Remove(tempFile)

	err := cfg.LoadFromFile(tempFile, opts)
	if err != nil {
		fmt.Printf("  Error loading with options: %v\n", err)
	} else {
		fmt.Println("  Advanced options loading successful:")
		fmt.Printf("    App Name: %s\n", cfg.GetString("app_name"))
		fmt.Printf("    Database URL: %s\n", cfg.GetString("database_url"))
		fmt.Printf("    Server Port (default applied): %d\n", cfg.GetInt("server_port"))
		fmt.Printf("    Debug Mode (default applied): %v\n", cfg.GetBool("debug_mode"))
		fmt.Printf("    Timeout (default applied): %v\n", cfg.GetDuration("timeout"))
	}

	// Demonstrate validation failure
	fmt.Println("  Demonstrating validation failure:")
	invalidConfig := map[string]any{
		"app_name":     "InvalidApp",
		"database_url": "postgres://localhost/myapp",
		"server_port":  80, // This will fail validation
	}

	tempFile2 := createTempConfig(invalidConfig)
	defer os.Remove(tempFile2)

	cfg2, _ := config.New()
	err = cfg2.LoadFromFile(tempFile2, opts)
	if err != nil {
		fmt.Printf("    Validation failed as expected: %v\n", err)
	}
	fmt.Println()
}

// Example 8: Runtime Configuration Changes
func demonstrateRuntimeChanges() {
	cfg, _ := config.New()

	// Initial configuration
	cfg.LoadFromMap(map[string]any{
		"feature_flags": map[string]any{
			"new_ui":    false,
			"beta_api":  true,
			"analytics": false,
		},
		"limits": map[string]any{
			"max_requests": 1000,
			"rate_limit":   100,
		},
	})

	fmt.Println("  Initial configuration:")
	fmt.Printf("    New UI: %v\n", cfg.GetBool("feature_flags.new_ui"))
	fmt.Printf("    Beta API: %v\n", cfg.GetBool("feature_flags.beta_api"))
	fmt.Printf("    Max Requests: %d\n", cfg.GetInt("limits.max_requests"))

	// Runtime changes
	fmt.Println("  Making runtime changes:")
	cfg.Set("feature_flags.new_ui", true)
	cfg.Set("limits.max_requests", 2000)
	cfg.Set("new_feature", "dynamic_value")

	fmt.Println("  After runtime changes:")
	fmt.Printf("    New UI: %v\n", cfg.GetBool("feature_flags.new_ui"))
	fmt.Printf("    Beta API: %v\n", cfg.GetBool("feature_flags.beta_api"))
	fmt.Printf("    Max Requests: %d\n", cfg.GetInt("limits.max_requests"))
	fmt.Printf("    New Feature: %s\n", cfg.GetString("new_feature"))

	// Demonstrate nested defaults
	cfg.SetNestedDefaults(map[string]any{
		"cache.ttl":      "1h",
		"cache.max_size": 1000,
		"existing_key":   "should_not_override", // This won't override existing value
	})

	fmt.Println("  After setting nested defaults:")
	fmt.Printf("    Cache TTL: %v\n", cfg.GetDuration("cache.ttl"))
	fmt.Printf("    Cache Max Size: %d\n", cfg.GetInt("cache.max_size"))
	fmt.Println()
}

// Example 9: Utility Methods
func demonstrateUtilityMethods() {
	cfg, _ := config.New()

	cfg.LoadFromMap(map[string]any{
		"app":     "UtilityDemo",
		"version": "1.0.0",
		"database": map[string]any{
			"host": "localhost",
			"port": 5432,
		},
		"features": []string{"auth", "api", "logging"},
	})

	fmt.Printf("  Configuration size: %d keys\n", cfg.Size())
	fmt.Printf("  Is empty: %v\n", cfg.IsEmpty())

	fmt.Println("  All keys:")
	for _, key := range cfg.Keys() {
		fmt.Printf("    - %s\n", key)
	}

	fmt.Println("  Checking key existence:")
	fmt.Printf("    Has 'app': %v\n", cfg.Has("app"))
	fmt.Printf("    Has 'database.host': %v\n", cfg.Has("database.host"))
	fmt.Printf("    Has 'nonexistent': %v\n", cfg.Has("nonexistent"))

	fmt.Println("  All configuration data:")
	allData := cfg.GetAll()
	for key, value := range allData {
		fmt.Printf("    %s: %v\n", key, value)
	}

	fmt.Println("  Configuration string representation:")
	configStr := cfg.String()
	fmt.Printf("    %s\n", configStr)

	// Clear and check
	fmt.Println("  After clearing:")
	cfg.Clear()
	fmt.Printf("    Size: %d\n", cfg.Size())
	fmt.Printf("    Is empty: %v\n", cfg.IsEmpty())
	fmt.Println()
}

// Example 10: Error Handling
func demonstrateErrorHandling() {
	fmt.Println("  Error handling scenarios:")

	// File not found
	cfg, _ := config.New()
	err := cfg.LoadFromFile("nonexistent.json", nil)
	if err != nil {
		fmt.Printf("    File not found error: %v\n", err)
	}

	// Invalid JSON
	invalidJSON := `{"invalid": json}`
	tempFile := filepath.Join(os.TempDir(), "invalid.json")
	os.WriteFile(tempFile, []byte(invalidJSON), 0644)
	defer os.Remove(tempFile)

	err = cfg.LoadFromFile(tempFile, nil)
	if err != nil {
		fmt.Printf("    Invalid JSON error: %v\n", err)
	}

	// Required keys missing
	opts := &config.LoadOptions{
		RequiredKeys: []string{"required_key"},
	}

	validJSON := `{"other_key": "value"}`
	tempFile2 := filepath.Join(os.TempDir(), "valid.json")
	os.WriteFile(tempFile2, []byte(validJSON), 0644)
	defer os.Remove(tempFile2)

	err = cfg.LoadFromFile(tempFile2, opts)
	if err != nil {
		fmt.Printf("    Required key missing error: %v\n", err)
	}

	// Nil config handling
	var nilCfg *config.Config
	fmt.Printf("    Nil config GetString: '%s'\n", nilCfg.GetString("key", "default"))
	fmt.Printf("    Nil config Has: %v\n", nilCfg.Has("key"))
	fmt.Printf("    Nil config Size: %d\n", nilCfg.Size())
	fmt.Println()
}

// Example 11: Thread Safety
func demonstrateThreadSafety() {
	cfg, _ := config.New()

	// Initial setup
	cfg.LoadFromMap(map[string]any{
		"counter": 0,
		"shared":  "initial",
	})

	fmt.Println("  Thread safety demonstration:")
	fmt.Println("    Starting concurrent operations...")

	// Simulate concurrent access
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 5; i++ {
			cfg.Set("counter", i)
			cfg.Set("shared", fmt.Sprintf("value_%d", i))
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 5; i++ {
			counter := cfg.GetInt("counter")
			shared := cfg.GetString("shared")
			_ = cfg.Has("counter")
			_ = cfg.Size()
			fmt.Printf("    Read - Counter: %d, Shared: %s\n", counter, shared)
			time.Sleep(15 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	fmt.Printf("    Final counter value: %d\n", cfg.GetInt("counter"))
	fmt.Printf("    Final shared value: %s\n", cfg.GetString("shared"))
	fmt.Println()
}

// Example 12: Real-World Application Configuration
func demonstrateRealWorldExample() {
	fmt.Println("  Real-world web application configuration:")

	// Create a comprehensive application configuration
	appConfig := map[string]any{
		"app": map[string]any{
			"name":        "MyWebApp",
			"version":     "1.2.3",
			"environment": "production",
			"debug":       false,
		},
		"server": map[string]any{
			"host":          "0.0.0.0",
			"port":          8080,
			"ssl_enabled":   true,
			"ssl_cert":      "/etc/ssl/certs/app.crt",
			"ssl_key":       "/etc/ssl/private/app.key",
			"read_timeout":  "30s",
			"write_timeout": "30s",
			"idle_timeout":  "120s",
		},
		"database": map[string]any{
			"driver":   "postgres",
			"host":     "db.example.com",
			"port":     5432,
			"name":     "webapp_prod",
			"user":     "webapp_user",
			"password": "secure_password",
			"ssl_mode": "require",
			"pool": map[string]any{
				"max_open":     25,
				"max_idle":     5,
				"max_lifetime": "5m",
			},
		},
		"redis": map[string]any{
			"host":      "redis.example.com",
			"port":      6379,
			"password":  "redis_password",
			"db":        0,
			"pool_size": 10,
		},
		"logging": map[string]any{
			"level":  "info",
			"format": "json",
			"output": "stdout",
			"file": map[string]any{
				"enabled":     false,
				"path":        "/var/log/webapp.log",
				"max_size":    "100MB",
				"max_backups": 3,
				"max_age":     30,
			},
		},
		"security": map[string]any{
			"jwt": map[string]any{
				"secret":      "super_secret_jwt_key",
				"expiration":  "24h",
				"refresh_exp": "168h",
			},
			"cors": map[string]any{
				"enabled":     true,
				"origins":     []string{"https://example.com", "https://app.example.com"},
				"methods":     []string{"GET", "POST", "PUT", "DELETE"},
				"headers":     []string{"Content-Type", "Authorization"},
				"credentials": true,
			},
			"rate_limit": map[string]any{
				"enabled":  true,
				"requests": 100,
				"window":   "1m",
				"burst":    10,
			},
		},
		"features": map[string]any{
			"user_registration":  true,
			"email_verification": true,
			"two_factor_auth":    false,
			"analytics":          true,
			"metrics":            true,
		},
		"external_services": map[string]any{
			"email": map[string]any{
				"provider": "sendgrid",
				"api_key":  "sendgrid_api_key",
				"from":     "noreply@example.com",
			},
			"s3": map[string]any{
				"bucket":     "webapp-uploads",
				"region":     "us-west-2",
				"access_key": "s3_access_key",
				"secret_key": "s3_secret_key",
			},
		},
	}

	cfg, _ := config.New()
	cfg.LoadFromMap(appConfig)

	// Demonstrate accessing various configuration sections
	fmt.Printf("    Application: %s v%s (%s)\n",
		cfg.GetString("app.name"),
		cfg.GetString("app.version"),
		cfg.GetString("app.environment"))

	fmt.Printf("    Server: %s:%d (SSL: %v)\n",
		cfg.GetString("server.host"),
		cfg.GetInt("server.port"),
		cfg.GetBool("server.ssl_enabled"))

	fmt.Printf("    Database: %s://%s:%d/%s (Pool: %d/%d)\n",
		cfg.GetString("database.driver"),
		cfg.GetString("database.host"),
		cfg.GetInt("database.port"),
		cfg.GetString("database.name"),
		cfg.GetInt("database.pool.max_idle"),
		cfg.GetInt("database.pool.max_open"))

	fmt.Printf("    Redis: %s:%d (Pool: %d)\n",
		cfg.GetString("redis.host"),
		cfg.GetInt("redis.port"),
		cfg.GetInt("redis.pool_size"))

	fmt.Printf("    Logging: %s level, %s format\n",
		cfg.GetString("logging.level"),
		cfg.GetString("logging.format"))

	fmt.Printf("    JWT Expiration: %v\n",
		cfg.GetDuration("security.jwt.expiration"))

	fmt.Printf("    CORS Origins: %v\n",
		cfg.GetStringSlice("security.cors.origins"))

	fmt.Printf("    Rate Limit: %d requests per %v\n",
		cfg.GetInt("security.rate_limit.requests"),
		cfg.GetDuration("security.rate_limit.window"))

	fmt.Printf("    Features Enabled: Registration=%v, 2FA=%v, Analytics=%v\n",
		cfg.GetBool("features.user_registration"),
		cfg.GetBool("features.two_factor_auth"),
		cfg.GetBool("features.analytics"))

	fmt.Printf("    Email Provider: %s\n",
		cfg.GetString("external_services.email.provider"))

	fmt.Printf("    S3 Bucket: %s in %s\n",
		cfg.GetString("external_services.s3.bucket"),
		cfg.GetString("external_services.s3.region"))

	// Demonstrate configuration validation for this real-world scenario
	fmt.Println("  Configuration validation:")

	requiredKeys := []string{
		"app.name", "app.version",
		"server.host", "server.port",
		"database.host", "database.name",
		"security.jwt.secret",
	}

	allPresent := true
	for _, key := range requiredKeys {
		if !cfg.Has(key) {
			fmt.Printf("    ❌ Missing required key: %s\n", key)
			allPresent = false
		}
	}

	if allPresent {
		fmt.Println("    ✅ All required configuration keys present")
	}

	fmt.Printf("    Total configuration keys: %d\n", cfg.Size())
	fmt.Println()
}

// Helper functions

func createSampleFiles() {
	// JSON configuration
	jsonConfig := `{
  "app_name": "JsonApp",
  "server_port": 8080,
  "debug_mode": true,
  "database": {
    "host": "localhost",
    "port": 5432,
    "name": "jsonapp"
  },
  "features": ["auth", "api", "web"],
  "timeout": "30s"
}`

	// YAML configuration
	yamlConfig := `app_name: YamlApp
server_port: 8081
debug_mode: false
database:
  host: yaml-db.example.com
  port: 5432
  name: yamlapp
features:
  - auth
  - api
  - dashboard
timeout: 45s`

	// INI configuration
	iniConfig := `app_name=IniApp
server_port=8082
debug_mode=true
timeout=60

[database]
host=ini-db.example.com
port=5432
name=iniapp

[features]
auth=true
api=true
reporting=false`

	// Write files
	os.WriteFile("config.json", []byte(jsonConfig), 0644)
	os.WriteFile("config.yaml", []byte(yamlConfig), 0644)
	os.WriteFile("config.ini", []byte(iniConfig), 0644)
}

func createTempConfig(data map[string]any) string {
	tempFile := filepath.Join(os.TempDir(), "temp_config.json")

	// Convert to JSON
	jsonData := fmt.Sprintf(`{
  "app_name": "%s",
  "database_url": "%s",
  "server_port": %v
}`, data["app_name"], data["database_url"], data["server_port"])

	os.WriteFile(tempFile, []byte(jsonData), 0644)
	return tempFile
}

func getMapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
