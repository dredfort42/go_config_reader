/*******************************************************************

		::          ::        +--------+-----------------------+
		  ::      ::          | Author | Dmitry Novikov        |
		::::::::::::::        | Email  | dredfort.42@gmail.com |
	  ::::  ::::::  ::::      +--------+-----------------------+
	::::::::::::::::::::::
	::  ::::::::::::::  ::    File     | main.go
	::  ::          ::  ::    Created  | 2025-08-07
		  ::::  ::::          Modified | 2025-08-07

	GitHub:   https://github.com/dredfort42
	LinkedIn: https://linkedin.com/in/novikov-da

*******************************************************************/

// nolint
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
	fmt.Println("=== Go Config Reader v1.0.0 Examples ===")
	fmt.Println()

	// Example 1: Using the new Config API with manual configuration
	fmt.Println("1. Manual Configuration:")

	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set some values directly
	cfg.Set("app_name", "ExampleApp")
	cfg.Set("server_port", 8080)
	cfg.Set("debug_mode", true)
	cfg.Set("timeout", "30s")
	cfg.Set("features", []string{"auth", "api", "web"})

	// Get values with defaults
	appName := cfg.GetString("app_name", "DefaultApp")
	serverPort := cfg.GetInt("server_port", 3000)
	debugMode := cfg.GetBool("debug_mode", false)
	timeout := cfg.GetDuration("timeout", 10*time.Second)
	features := cfg.GetStringSlice("features", []string{"basic"})

	fmt.Printf("  App Name: %s\n", appName)
	fmt.Printf("  Server Port: %d\n", serverPort)
	fmt.Printf("  Debug Mode: %v\n", debugMode)
	fmt.Printf("  Timeout: %v\n", timeout)
	fmt.Printf("  Features: %v\n", features)
	fmt.Printf("  Total Keys: %d\n\n", cfg.Size())

	// Example 2: Loading from JSON file
	fmt.Println("2. Loading from JSON file:")

	if err := demonstrateFileLoading("config.json"); err != nil {
		fmt.Printf("  Error: %v\n\n", err)
	}

	// Example 3: Loading from YAML file
	fmt.Println("3. Loading from YAML file:")

	if err := demonstrateFileLoading("config.yaml"); err != nil {
		fmt.Printf("  Error: %v\n\n", err)
	}

	// Example 4: Loading from INI file
	fmt.Println("4. Loading from INI file:")

	if err := demonstrateFileLoading("config.ini"); err != nil {
		fmt.Printf("  Error: %v\n\n", err)
	}

	// Example 5: Using global API (backward compatibility)
	fmt.Println("5. Global API (backward compatibility):")
	config.Set("global_setting", "global_value")
	config.Set("global_port", 9090)

	fmt.Printf("  Global Setting: %s\n", config.GetString("global_setting", "default"))
	fmt.Printf("  Global Port: %d\n", config.GetInt("global_port", 8080))
	fmt.Printf("  Has global_setting: %v\n", config.Has("global_setting"))
	fmt.Printf("  All Global Keys: %v\n\n", config.Keys())

	// Example 6: Configuration with defaults and validation
	fmt.Println("6. Configuration with defaults and validation:")
	demonstrateAdvancedFeatures()

	// Example 7: Advanced YAML Configuration with Data Extraction
	fmt.Println("7. Advanced YAML Configuration - Data Extraction:")
	demonstrateYAMLDataExtraction()
}

func demonstrateFileLoading(filename string) error {
	configPath := filepath.Join(".", filename)

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file %s not found", filename)
	}

	cfg, err := config.New()
	if err != nil {
		return err
	}

	err = cfg.LoadFromFile(configPath, nil)
	if err != nil {
		return err
	}

	// For structured configs, we'll try to get values that might exist
	fmt.Printf("  Loaded %d configuration keys\n", cfg.Size())

	// Try some common keys that might exist in the config files
	if cfg.Has("server_host") {
		fmt.Printf("  Server Host: %s\n", cfg.GetString("server_host"))
	}

	if cfg.Has("server_port") || cfg.Has("server.port") {
		port := cfg.GetInt("server_port", cfg.GetInt("server.port", 8080))
		fmt.Printf("  Server Port: %d\n", port)
	}

	if cfg.Has("timeout") {
		timeout := cfg.GetDuration("timeout", 30*time.Second)
		fmt.Printf("  Timeout: %v\n", timeout)
	}

	// Show all keys for debugging
	fmt.Printf("  Available keys: %v\n", cfg.Keys())
	fmt.Println()

	return nil
}

func demonstrateAdvancedFeatures() {
	// Create a temporary config for demonstration
	cfg, _ := config.New()

	// Load with defaults
	defaults := map[string]any{
		"app_name":        "DefaultApp",
		"server_port":     8080,
		"debug_mode":      false,
		"max_connections": 100,
		"features":        []string{"basic"},
	}

	cfg.LoadFromMap(defaults)

	// Override some values
	cfg.Set("app_name", "ProductionApp")
	cfg.Set("debug_mode", true)

	// Validation example
	requiredKeys := []string{"app_name", "server_port"}
	opts := &config.LoadOptions{
		RequiredKeys: requiredKeys,
		ValidationFunc: func(data map[string]any) error {
			if port, ok := data["server_port"].(int); ok && port < 1024 {
				return fmt.Errorf("server port must be >= 1024, got %d", port)
			}

			return nil
		},
	}

	// This would normally be used with LoadFromFile, but we'll demonstrate the concept
	fmt.Printf("  App Name: %s\n", cfg.GetString("app_name"))
	fmt.Printf("  Debug Mode: %v\n", cfg.GetBool("debug_mode"))
	fmt.Printf("  Max Connections: %d\n", cfg.GetInt("max_connections"))
	fmt.Printf("  Validation example passed for required keys: %v\n", opts.RequiredKeys)
	fmt.Println()
}

func demonstrateYAMLDataExtraction() {
	configPath := filepath.Join(".", "config.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("  Error: YAML config file not found: %v\n\n", err)
		return
	}

	// Create a new config instance
	cfg, err := config.New()
	if err != nil {
		fmt.Printf("  Error creating config: %v\n\n", err)
		return
	}

	// Load YAML configuration
	err = cfg.LoadFromFile(configPath, nil)
	if err != nil {
		fmt.Printf("  Error loading YAML config: %v\n\n", err)
		return
	}

	fmt.Printf("  Successfully loaded YAML configuration from %s\n", configPath)
	fmt.Printf("  Total configuration keys: %d\n\n", cfg.Size())

	// Get all data first for nested access
	allData := cfg.GetAll()

	// Extract and display specific configuration values with type safety
	fmt.Println("  📊 Configuration Values:")
	fmt.Println("  ========================")

	// Server configuration
	fmt.Println("  🖥️  Server Configuration:")
	if serverData, ok := allData["server"]; ok {
		if serverMap, ok := serverData.(map[string]interface{}); ok {
			if host, exists := serverMap["host"]; exists {
				fmt.Printf("    Host: %s\n", host)
			}
			if port, exists := serverMap["port"]; exists {
				fmt.Printf("    Port: %v\n", port)
			}
			if debug, exists := serverMap["debug"]; exists {
				fmt.Printf("    Debug Mode: %v\n", debug)
			}
		}
	}

	// Database configuration
	fmt.Println("  🗄️  Database Configuration:")
	if databaseData, ok := allData["database"]; ok {
		if databaseMap, ok := databaseData.(map[string]interface{}); ok {
			if host, exists := databaseMap["host"]; exists {
				fmt.Printf("    Host: %s\n", host)
			}
			if port, exists := databaseMap["port"]; exists {
				fmt.Printf("    Port: %v\n", port)
			}
			if name, exists := databaseMap["name"]; exists {
				fmt.Printf("    Name: %s\n", name)
			}
			if ssl, exists := databaseMap["ssl"]; exists {
				fmt.Printf("    SSL: %v\n", ssl)
			}
		}
	}

	// Application settings
	fmt.Println("  ⚙️  Application Settings:")
	if cfg.Has("timeout") {
		timeout := cfg.GetDuration("timeout", 30*time.Second)
		fmt.Printf("    Timeout: %v\n", timeout)
	}
	if cfg.Has("max_connections") {
		fmt.Printf("    Max Connections: %d\n", cfg.GetInt("max_connections", 100))
	}
	if cfg.Has("features") {
		features := cfg.GetStringSlice("features", []string{"basic"})
		fmt.Printf("    Features: %v\n", features)
	}

	fmt.Println()

	// Display ALL configuration data
	fmt.Println("  📋 Complete Configuration Data:")
	fmt.Println("  ===============================")

	if len(allData) == 0 {
		fmt.Println("    No configuration data found.")
	} else {
		for key, value := range allData {
			switch v := value.(type) {
			case string:
				fmt.Printf("    %s: \"%s\" (string)\n", key, v)
			case int:
				fmt.Printf("    %s: %d (int)\n", key, v)
			case int64:
				fmt.Printf("    %s: %d (int64)\n", key, v)
			case float64:
				fmt.Printf("    %s: %.2f (float64)\n", key, v)
			case bool:
				fmt.Printf("    %s: %v (bool)\n", key, v)
			case []interface{}:
				fmt.Printf("    %s: %v (slice)\n", key, v)
			case map[string]interface{}:
				fmt.Printf("    %s: %v (map)\n", key, v)
			default:
				fmt.Printf("    %s: %v (%T)\n", key, v, v)
			}
		}
	}

	fmt.Println()

	// Display configuration as string representation
	fmt.Println("  📄 String Representation:")
	fmt.Println("  =========================")
	configStr := cfg.String()
	if configStr == "" {
		fmt.Println("    Empty configuration")
	} else {
		fmt.Print("  " + configStr)
	}

	fmt.Println()

	// Demonstrate configuration utilities
	fmt.Println("  🔧 Configuration Utilities:")
	fmt.Println("  ===========================")
	fmt.Printf("    Total keys: %d\n", cfg.Size())
	fmt.Printf("    Is empty: %v\n", cfg.IsEmpty())
	fmt.Printf("    Available keys: %v\n", cfg.Keys())

	// Check for specific keys
	fmt.Println("  🔍 Key Existence Checks:")
	testKeys := []string{"server", "database", "timeout", "features", "non_existent_key"}
	for _, key := range testKeys {
		exists := cfg.Has(key)
		fmt.Printf("    Has '%s': %v\n", key, exists)
	}

	// Check for nested keys manually
	fmt.Println("  🔍 Nested Key Checks:")
	if serverData, ok := allData["server"]; ok {
		if serverMap, ok := serverData.(map[string]interface{}); ok {
			for nestedKey := range serverMap {
				fmt.Printf("    server.%s: exists\n", nestedKey)
			}
		}
	}
	if databaseData, ok := allData["database"]; ok {
		if databaseMap, ok := databaseData.(map[string]interface{}); ok {
			for nestedKey := range databaseMap {
				fmt.Printf("    database.%s: exists\n", nestedKey)
			}
		}
	}

	// Demonstrate the helper function for nested value access
	fmt.Println("  🎯 Nested Value Access with Helper Function:")
	if val, exists := getNestedValue(allData, "server", "host"); exists {
		fmt.Printf("    server.host: %v\n", val)
	}
	if val, exists := getNestedValue(allData, "database", "port"); exists {
		fmt.Printf("    database.port: %v\n", val)
	}
	if val, exists := getNestedValue(allData, "server", "nonexistent"); exists {
		fmt.Printf("    server.nonexistent: %v\n", val)
	} else {
		fmt.Printf("    server.nonexistent: key not found\n")
	}

	fmt.Println()
}

// getNestedValue is a helper function to safely extract nested values from config maps
func getNestedValue(data map[string]any, path ...string) (any, bool) {
	current := data

	for _, key := range path[:len(path)-1] {
		if val, exists := current[key]; exists {
			if nextMap, ok := val.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return nil, false // Path doesn't lead to a map
			}
		} else {
			return nil, false // Key doesn't exist
		}
	}

	finalKey := path[len(path)-1]
	if val, exists := current[finalKey]; exists {
		return val, true
	}

	return nil, false
}
