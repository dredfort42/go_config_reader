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

// Example: Microservice Configuration Management
package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	config "github.com/dredfort42/go_config_reader"
)

// ServiceConfig represents microservice configuration
type ServiceConfig struct {
	Name      string
	Version   string
	Port      int
	Health    HealthConfig
	Metrics   MetricsConfig
	Tracing   TracingConfig
	Discovery ServiceDiscoveryConfig
	Circuit   CircuitBreakerConfig
}

type HealthConfig struct {
	Enabled  bool
	Endpoint string
	Timeout  time.Duration
}

type MetricsConfig struct {
	Enabled  bool
	Endpoint string
	Interval time.Duration
}

type TracingConfig struct {
	Enabled     bool
	ServiceName string
	Endpoint    string
	SampleRate  float64
}

type ServiceDiscoveryConfig struct {
	Enabled bool
	Type    string // consul, etcd, kubernetes
	Address string
	TTL     time.Duration
}

type CircuitBreakerConfig struct {
	Enabled          bool
	MaxRequests      int
	Interval         time.Duration
	Timeout          time.Duration
	FailureThreshold int
	SuccessThreshold int
}

func main() {
	fmt.Println("🔬 Microservice Configuration Example")
	fmt.Println("======================================\n")

	// Demonstrate different environment configurations
	environments := []string{"development", "staging", "production"}

	for _, env := range environments {
		fmt.Printf("🌍 Environment: %s\n", env)
		fmt.Println(strings.Repeat("-", 40))

		cfg := loadMicroserviceConfig(env)
		serviceConfig := parseMicroserviceConfig(cfg)
		displayMicroserviceConfig(serviceConfig)

		fmt.Println()
	}

	// Demonstrate configuration hot-reloading simulation
	fmt.Println("🔄 Configuration Hot-Reload Simulation")
	fmt.Println("=======================================")
	demonstrateHotReload()
}

func loadMicroserviceConfig(environment string) *config.Config {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to create config: %v", err)
	}

	// Environment-specific defaults
	defaults := getEnvironmentDefaults(environment)

	// Advanced options with environment-specific validation
	opts := &config.LoadOptions{
		DefaultValues: defaults,
		RequiredKeys: []string{
			"service.name",
			"service.version",
			"service.port",
		},
		ValidationFunc: func(data map[string]any) error {
			return validateMicroserviceConfig(data, environment)
		},
		IgnoreEnv: false, // Allow environment variable overrides
	}

	// Try environment-specific config files
	configFiles := []string{
		fmt.Sprintf("microservice.%s.json", environment),
		fmt.Sprintf("microservice.%s.yaml", environment),
		"microservice.json",
		"microservice.yaml",
	}

	loaded := false
	for _, file := range configFiles {
		if err := cfg.LoadFromFile(file, opts); err == nil {
			fmt.Printf("  ✅ Loaded from: %s\n", file)
			loaded = true
			break
		}
	}

	if !loaded {
		fmt.Printf("  ⚠️  No config file found, using %s defaults\n", environment)
		cfg.LoadFromMap(defaults)
	}

	return cfg
}

func getEnvironmentDefaults(environment string) map[string]any {
	baseDefaults := map[string]any{
		"service.name":    "example-microservice",
		"service.version": "1.0.0",
		"service.port":    8080,

		"health.enabled":  true,
		"health.endpoint": "/health",
		"health.timeout":  "5s",

		"metrics.enabled":  true,
		"metrics.endpoint": "/metrics",
		"metrics.interval": "30s",

		"tracing.enabled":      false,
		"tracing.service_name": "example-microservice",
		"tracing.endpoint":     "http://jaeger:14268/api/traces",
		"tracing.sample_rate":  0.1,

		"discovery.enabled": false,
		"discovery.type":    "consul",
		"discovery.address": "consul:8500",
		"discovery.ttl":     "30s",

		"circuit.enabled":           false,
		"circuit.max_requests":      100,
		"circuit.interval":          "10s",
		"circuit.timeout":           "30s",
		"circuit.failure_threshold": 5,
		"circuit.success_threshold": 3,
	}

	// Environment-specific overrides
	switch environment {
	case "development":
		baseDefaults["service.port"] = 8080
		baseDefaults["tracing.enabled"] = true
		baseDefaults["tracing.sample_rate"] = 1.0 // Sample everything in dev
		baseDefaults["metrics.interval"] = "10s"  // More frequent metrics in dev

	case "staging":
		baseDefaults["service.port"] = 8081
		baseDefaults["tracing.enabled"] = true
		baseDefaults["tracing.sample_rate"] = 0.5
		baseDefaults["discovery.enabled"] = true
		baseDefaults["circuit.enabled"] = true

	case "production":
		baseDefaults["service.port"] = 8080
		baseDefaults["tracing.enabled"] = true
		baseDefaults["tracing.sample_rate"] = 0.1
		baseDefaults["discovery.enabled"] = true
		baseDefaults["circuit.enabled"] = true
		baseDefaults["circuit.max_requests"] = 1000
		baseDefaults["metrics.interval"] = "60s"
	}

	return baseDefaults
}

func validateMicroserviceConfig(data map[string]any, environment string) error {
	// Port validation
	if port, ok := data["service.port"].(int); ok {
		if port < 1024 || port > 65535 {
			return fmt.Errorf("service port must be between 1024 and 65535, got %d", port)
		}
	}

	// Tracing sample rate validation
	if rate, ok := data["tracing.sample_rate"].(float64); ok {
		if rate < 0.0 || rate > 1.0 {
			return fmt.Errorf("tracing sample rate must be between 0.0 and 1.0, got %f", rate)
		}
	}

	// Circuit breaker validation
	if enabled, ok := data["circuit.enabled"].(bool); ok && enabled {
		if threshold, ok := data["circuit.failure_threshold"].(int); ok {
			if threshold < 1 {
				return fmt.Errorf("circuit breaker failure threshold must be at least 1, got %d", threshold)
			}
		}
	}

	// Production-specific validations
	if environment == "production" {
		if enabled, ok := data["discovery.enabled"].(bool); ok && !enabled {
			return fmt.Errorf("service discovery must be enabled in production")
		}
		if enabled, ok := data["circuit.enabled"].(bool); ok && !enabled {
			return fmt.Errorf("circuit breaker must be enabled in production")
		}
	}

	return nil
}

func parseMicroserviceConfig(cfg *config.Config) ServiceConfig {
	return ServiceConfig{
		Name:    cfg.GetString("service.name"),
		Version: cfg.GetString("service.version"),
		Port:    cfg.GetInt("service.port"),

		Health: HealthConfig{
			Enabled:  cfg.GetBool("health.enabled"),
			Endpoint: cfg.GetString("health.endpoint"),
			Timeout:  cfg.GetDuration("health.timeout"),
		},

		Metrics: MetricsConfig{
			Enabled:  cfg.GetBool("metrics.enabled"),
			Endpoint: cfg.GetString("metrics.endpoint"),
			Interval: cfg.GetDuration("metrics.interval"),
		},

		Tracing: TracingConfig{
			Enabled:     cfg.GetBool("tracing.enabled"),
			ServiceName: cfg.GetString("tracing.service_name"),
			Endpoint:    cfg.GetString("tracing.endpoint"),
			SampleRate:  cfg.GetFloat64("tracing.sample_rate"),
		},

		Discovery: ServiceDiscoveryConfig{
			Enabled: cfg.GetBool("discovery.enabled"),
			Type:    cfg.GetString("discovery.type"),
			Address: cfg.GetString("discovery.address"),
			TTL:     cfg.GetDuration("discovery.ttl"),
		},

		Circuit: CircuitBreakerConfig{
			Enabled:          cfg.GetBool("circuit.enabled"),
			MaxRequests:      cfg.GetInt("circuit.max_requests"),
			Interval:         cfg.GetDuration("circuit.interval"),
			Timeout:          cfg.GetDuration("circuit.timeout"),
			FailureThreshold: cfg.GetInt("circuit.failure_threshold"),
			SuccessThreshold: cfg.GetInt("circuit.success_threshold"),
		},
	}
}

func displayMicroserviceConfig(sc ServiceConfig) {
	fmt.Printf("  📦 Service: %s v%s (Port: %d)\n", sc.Name, sc.Version, sc.Port)

	fmt.Printf("  ❤️  Health Check: %v", sc.Health.Enabled)
	if sc.Health.Enabled {
		fmt.Printf(" (%s, timeout: %v)", sc.Health.Endpoint, sc.Health.Timeout)
	}
	fmt.Println()

	fmt.Printf("  📊 Metrics: %v", sc.Metrics.Enabled)
	if sc.Metrics.Enabled {
		fmt.Printf(" (%s, interval: %v)", sc.Metrics.Endpoint, sc.Metrics.Interval)
	}
	fmt.Println()

	fmt.Printf("  🔍 Tracing: %v", sc.Tracing.Enabled)
	if sc.Tracing.Enabled {
		fmt.Printf(" (rate: %.1f, endpoint: %s)", sc.Tracing.SampleRate, sc.Tracing.Endpoint)
	}
	fmt.Println()

	fmt.Printf("  🌐 Service Discovery: %v", sc.Discovery.Enabled)
	if sc.Discovery.Enabled {
		fmt.Printf(" (%s at %s, TTL: %v)", sc.Discovery.Type, sc.Discovery.Address, sc.Discovery.TTL)
	}
	fmt.Println()

	fmt.Printf("  ⚡ Circuit Breaker: %v", sc.Circuit.Enabled)
	if sc.Circuit.Enabled {
		fmt.Printf(" (failure threshold: %d, max requests: %d)",
			sc.Circuit.FailureThreshold, sc.Circuit.MaxRequests)
	}
	fmt.Println()
}

func demonstrateHotReload() {
	cfg, _ := config.New()

	// Initial configuration
	initialConfig := map[string]any{
		"service.name":         "hot-reload-service",
		"service.version":      "1.0.0",
		"circuit.enabled":      false,
		"circuit.max_requests": 100,
		"metrics.interval":     "30s",
		"tracing.sample_rate":  0.1,
	}

	cfg.LoadFromMap(initialConfig)

	fmt.Println("  Initial Configuration:")
	fmt.Printf("    Circuit Breaker: %v\n", cfg.GetBool("circuit.enabled"))
	fmt.Printf("    Max Requests: %d\n", cfg.GetInt("circuit.max_requests"))
	fmt.Printf("    Metrics Interval: %v\n", cfg.GetDuration("metrics.interval"))
	fmt.Printf("    Tracing Sample Rate: %.1f\n", cfg.GetFloat64("tracing.sample_rate"))

	// Simulate configuration updates (hot reload)
	fmt.Println("\n  Simulating hot reload with updates...")

	// Update 1: Enable circuit breaker
	cfg.Set("circuit.enabled", true)
	cfg.Set("circuit.max_requests", 500)
	fmt.Println("    ✅ Enabled circuit breaker with 500 max requests")

	// Update 2: Increase metrics frequency
	cfg.Set("metrics.interval", "10s")
	fmt.Println("    ✅ Increased metrics frequency to 10s")

	// Update 3: Increase tracing sample rate
	cfg.Set("tracing.sample_rate", 0.5)
	fmt.Println("    ✅ Increased tracing sample rate to 50%")

	// Show updated configuration
	fmt.Println("\n  Updated Configuration:")
	fmt.Printf("    Circuit Breaker: %v\n", cfg.GetBool("circuit.enabled"))
	fmt.Printf("    Max Requests: %d\n", cfg.GetInt("circuit.max_requests"))
	fmt.Printf("    Metrics Interval: %v\n", cfg.GetDuration("metrics.interval"))
	fmt.Printf("    Tracing Sample Rate: %.1f\n", cfg.GetFloat64("tracing.sample_rate"))

	// Demonstrate configuration validation during hot reload
	fmt.Println("\n  Testing validation during hot reload:")

	// Try invalid sample rate
	fmt.Println("    Attempting to set invalid tracing sample rate (1.5)...")
	originalRate := cfg.GetFloat64("tracing.sample_rate")
	cfg.Set("tracing.sample_rate", 1.5)

	// Validate the configuration
	if rate := cfg.GetFloat64("tracing.sample_rate"); rate > 1.0 {
		fmt.Printf("    ❌ Invalid sample rate detected: %.1f, reverting to %.1f\n", rate, originalRate)
		cfg.Set("tracing.sample_rate", originalRate)
	}

	// Show configuration keys and size
	fmt.Printf("\n  Configuration Summary:\n")
	fmt.Printf("    Total Keys: %d\n", cfg.Size())
	fmt.Printf("    Available Keys: %v\n", cfg.Keys())
}

// Helper function since strings is not imported
func stringRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// Create sample environment-specific configuration files
func createSampleMicroserviceConfigs() {
	// Development configuration
	devConfig := `{
  "service": {
    "name": "example-microservice",
    "version": "1.0.0",
    "port": 8080
  },
  "health": {
    "enabled": true,
    "endpoint": "/health",
    "timeout": "5s"
  },
  "metrics": {
    "enabled": true,
    "endpoint": "/metrics",
    "interval": "10s"
  },
  "tracing": {
    "enabled": true,
    "service_name": "example-microservice",
    "endpoint": "http://localhost:14268/api/traces",
    "sample_rate": 1.0
  },
  "discovery": {
    "enabled": false
  },
  "circuit": {
    "enabled": false
  }
}`

	// Production configuration
	prodConfig := `{
  "service": {
    "name": "example-microservice",
    "version": "1.0.0",
    "port": 8080
  },
  "health": {
    "enabled": true,
    "endpoint": "/health",
    "timeout": "5s"
  },
  "metrics": {
    "enabled": true,
    "endpoint": "/metrics",
    "interval": "60s"
  },
  "tracing": {
    "enabled": true,
    "service_name": "example-microservice",
    "endpoint": "http://jaeger:14268/api/traces",
    "sample_rate": 0.1
  },
  "discovery": {
    "enabled": true,
    "type": "consul",
    "address": "consul:8500",
    "ttl": "30s"
  },
  "circuit": {
    "enabled": true,
    "max_requests": 1000,
    "interval": "10s",
    "timeout": "30s",
    "failure_threshold": 5,
    "success_threshold": 3
  }
}`

	// Write sample configuration files
	os.WriteFile("microservice.development.json", []byte(devConfig), 0644)
	os.WriteFile("microservice.production.json", []byte(prodConfig), 0644)

	fmt.Println("💡 Sample configuration files created:")
	fmt.Println("   - microservice.development.json")
	fmt.Println("   - microservice.production.json")
	fmt.Println("\nYou can now test with these files!")
}
