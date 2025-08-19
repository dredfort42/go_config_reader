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

// Example: HTTP Server Configuration with go_config_reader
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	config "github.com/dredfort42/go_config_reader"
)

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host           string
	Port           int
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	MaxHeaderBytes int
	SSLEnabled     bool
	SSLCertFile    string
	SSLKeyFile     string
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver      string
	Host        string
	Port        int
	Name        string
	User        string
	Password    string
	SSLMode     string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
}

// AppConfig represents complete application configuration
type AppConfig struct {
	AppName     string
	Version     string
	Environment string
	Debug       bool
	LogLevel    string
	Server      ServerConfig
	Database    DatabaseConfig
}

func main() {
	fmt.Println("🌐 HTTP Server Configuration Example")
	fmt.Println("=====================================\n")

	// Load configuration from file or use defaults
	cfg := loadConfiguration()

	// Parse configuration into structured types
	appConfig := parseConfiguration(cfg)

	// Display configuration
	displayConfiguration(appConfig)

	// Demonstrate server creation with configuration
	createServer(appConfig)
}

func loadConfiguration() *config.Config {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to create config instance: %v", err)
	}

	// Define comprehensive defaults
	defaults := map[string]any{
		"app.name":        "WebServer",
		"app.version":     "1.0.0",
		"app.environment": "development",
		"app.debug":       true,
		"app.log_level":   "info",

		"server.host":             "localhost",
		"server.port":             8080,
		"server.read_timeout":     "30s",
		"server.write_timeout":    "30s",
		"server.idle_timeout":     "120s",
		"server.max_header_bytes": 1048576, // 1MB
		"server.ssl_enabled":      false,
		"server.ssl_cert_file":    "",
		"server.ssl_key_file":     "",

		"database.driver":        "postgres",
		"database.host":          "localhost",
		"database.port":          5432,
		"database.name":          "webserver_db",
		"database.user":          "webserver_user",
		"database.password":      "secure_password",
		"database.ssl_mode":      "disable",
		"database.max_open_conn": 25,
		"database.max_idle_conn": 5,
		"database.max_lifetime":  "5m",
	}

	// Try to load from config file with defaults and validation
	opts := &config.LoadOptions{
		DefaultValues: defaults,
		RequiredKeys: []string{
			"app.name",
			"server.host",
			"server.port",
			"database.host",
			"database.name",
		},
		ValidationFunc: func(data map[string]any) error {
			// Validate port range
			if port, ok := data["server.port"].(int); ok {
				if port < 1 || port > 65535 {
					return fmt.Errorf("server port must be between 1 and 65535, got %d", port)
				}
			}

			// Validate database port
			if dbPort, ok := data["database.port"].(int); ok {
				if dbPort < 1 || dbPort > 65535 {
					return fmt.Errorf("database port must be between 1 and 65535, got %d", dbPort)
				}
			}

			// Validate log level
			if logLevel, ok := data["app.log_level"].(string); ok {
				validLevels := []string{"debug", "info", "warn", "error"}
				valid := false
				for _, level := range validLevels {
					if logLevel == level {
						valid = true
						break
					}
				}
				if !valid {
					return fmt.Errorf("invalid log level: %s, must be one of %v", logLevel, validLevels)
				}
			}

			return nil
		},
	}

	// Try to load from file, fall back to defaults if file doesn't exist
	configFiles := []string{"server_config.json", "server_config.yaml", "config.json", "config.yaml"}
	loaded := false

	for _, file := range configFiles {
		err := cfg.LoadFromFile(file, opts)
		if err == nil {
			fmt.Printf("✅ Loaded configuration from: %s\n", file)
			loaded = true
			break
		}
	}

	if !loaded {
		fmt.Println("⚠️  No configuration file found, using defaults")
		// Apply defaults manually since no file was loaded
		cfg.LoadFromMap(defaults)
	}

	return cfg
}

func parseConfiguration(cfg *config.Config) AppConfig {
	return AppConfig{
		AppName:     cfg.GetString("app.name"),
		Version:     cfg.GetString("app.version"),
		Environment: cfg.GetString("app.environment"),
		Debug:       cfg.GetBool("app.debug"),
		LogLevel:    cfg.GetString("app.log_level"),

		Server: ServerConfig{
			Host:           cfg.GetString("server.host"),
			Port:           cfg.GetInt("server.port"),
			ReadTimeout:    cfg.GetDuration("server.read_timeout"),
			WriteTimeout:   cfg.GetDuration("server.write_timeout"),
			IdleTimeout:    cfg.GetDuration("server.idle_timeout"),
			MaxHeaderBytes: cfg.GetInt("server.max_header_bytes"),
			SSLEnabled:     cfg.GetBool("server.ssl_enabled"),
			SSLCertFile:    cfg.GetString("server.ssl_cert_file"),
			SSLKeyFile:     cfg.GetString("server.ssl_key_file"),
		},

		Database: DatabaseConfig{
			Driver:      cfg.GetString("database.driver"),
			Host:        cfg.GetString("database.host"),
			Port:        cfg.GetInt("database.port"),
			Name:        cfg.GetString("database.name"),
			User:        cfg.GetString("database.user"),
			Password:    cfg.GetString("database.password"),
			SSLMode:     cfg.GetString("database.ssl_mode"),
			MaxOpenConn: cfg.GetInt("database.max_open_conn"),
			MaxIdleConn: cfg.GetInt("database.max_idle_conn"),
			MaxLifetime: cfg.GetDuration("database.max_lifetime"),
		},
	}
}

func displayConfiguration(appConfig AppConfig) {
	fmt.Printf("📋 Application Configuration:\n")
	fmt.Printf("   Name: %s v%s\n", appConfig.AppName, appConfig.Version)
	fmt.Printf("   Environment: %s\n", appConfig.Environment)
	fmt.Printf("   Debug Mode: %v\n", appConfig.Debug)
	fmt.Printf("   Log Level: %s\n", appConfig.LogLevel)
	fmt.Println()

	fmt.Printf("🌐 Server Configuration:\n")
	fmt.Printf("   Address: %s:%d\n", appConfig.Server.Host, appConfig.Server.Port)
	fmt.Printf("   Read Timeout: %v\n", appConfig.Server.ReadTimeout)
	fmt.Printf("   Write Timeout: %v\n", appConfig.Server.WriteTimeout)
	fmt.Printf("   Idle Timeout: %v\n", appConfig.Server.IdleTimeout)
	fmt.Printf("   Max Header Bytes: %d\n", appConfig.Server.MaxHeaderBytes)
	fmt.Printf("   SSL Enabled: %v\n", appConfig.Server.SSLEnabled)
	if appConfig.Server.SSLEnabled {
		fmt.Printf("   SSL Cert: %s\n", appConfig.Server.SSLCertFile)
		fmt.Printf("   SSL Key: %s\n", appConfig.Server.SSLKeyFile)
	}
	fmt.Println()

	fmt.Printf("🗄️  Database Configuration:\n")
	fmt.Printf("   Driver: %s\n", appConfig.Database.Driver)
	fmt.Printf("   Address: %s:%d\n", appConfig.Database.Host, appConfig.Database.Port)
	fmt.Printf("   Database: %s\n", appConfig.Database.Name)
	fmt.Printf("   User: %s\n", appConfig.Database.User)
	fmt.Printf("   SSL Mode: %s\n", appConfig.Database.SSLMode)
	fmt.Printf("   Connection Pool: %d/%d (max idle/open)\n", appConfig.Database.MaxIdleConn, appConfig.Database.MaxOpenConn)
	fmt.Printf("   Max Lifetime: %v\n", appConfig.Database.MaxLifetime)
	fmt.Println()
}

func createServer(appConfig AppConfig) {
	// Create HTTP server with configuration
	mux := http.NewServeMux()

	// Add some example routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := fmt.Sprintf(`{
  "app": "%s",
  "version": "%s",
  "environment": "%s",
  "debug": %v,
  "timestamp": "%s"
}`, appConfig.AppName, appConfig.Version, appConfig.Environment, appConfig.Debug, time.Now().UTC().Format(time.RFC3339))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "database": "connected"}`))
	})

	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		// Return non-sensitive configuration
		response := fmt.Sprintf(`{
  "server": {
    "host": "%s",
    "port": %d,
    "ssl_enabled": %v
  },
  "database": {
    "driver": "%s",
    "host": "%s",
    "port": %d,
    "name": "%s",
    "pool": {
      "max_open": %d,
      "max_idle": %d
    }
  }
}`, appConfig.Server.Host, appConfig.Server.Port, appConfig.Server.SSLEnabled,
			appConfig.Database.Driver, appConfig.Database.Host, appConfig.Database.Port,
			appConfig.Database.Name, appConfig.Database.MaxOpenConn, appConfig.Database.MaxIdleConn)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})

	// Create server instance
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", appConfig.Server.Host, appConfig.Server.Port),
		Handler:        mux,
		ReadTimeout:    appConfig.Server.ReadTimeout,
		WriteTimeout:   appConfig.Server.WriteTimeout,
		IdleTimeout:    appConfig.Server.IdleTimeout,
		MaxHeaderBytes: appConfig.Server.MaxHeaderBytes,
	}

	fmt.Printf("🚀 Server Configuration Complete\n")
	fmt.Printf("   Server would start on: %s\n", server.Addr)
	fmt.Printf("   Available endpoints:\n")
	fmt.Printf("     GET /        - Application info\n")
	fmt.Printf("     GET /health  - Health check\n")
	fmt.Printf("     GET /config  - Configuration info\n")

	if appConfig.Server.SSLEnabled {
		fmt.Printf("   SSL Configuration:\n")
		fmt.Printf("     Cert File: %s\n", appConfig.Server.SSLCertFile)
		fmt.Printf("     Key File: %s\n", appConfig.Server.SSLKeyFile)
		fmt.Println("   Would start with: server.ListenAndServeTLS()")
	} else {
		fmt.Println("   Would start with: server.ListenAndServe()")
	}

	fmt.Println("\n💡 To create a real server configuration file, create one of:")
	fmt.Println("   - server_config.json")
	fmt.Println("   - server_config.yaml")
	fmt.Println("   - config.json")
	fmt.Println("   - config.yaml")
	fmt.Println("\n   Environment variables will override file values for any matching keys.")
}
