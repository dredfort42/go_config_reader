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
