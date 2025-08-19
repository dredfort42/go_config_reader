# Go Configuration Reader

A modern, flexible, and thread-safe configuration library for Go applications that supports multiple formats and provides a clean API.

## Features

-   **Multiple Format Support**: JSON, YAML, and INI files
-   **Environment Variable Override**: Automatic environment variable priority
-   **Thread-Safe**: Concurrent access protection with RWMutex
-   **Type Safety**: Strong typing with automatic type conversion
-   **Default Values**: Built-in support for default values
-   **Nested Configuration**: Dot notation access for nested structures
-   **Validation**: Custom validation functions and required key checking
-   **Instance-Based API**: Clean instance-based configuration management
-   **Minimal Dependencies**: Only requires `gopkg.in/yaml.v3` for YAML support
-   **High Performance**: Optimized for speed with comprehensive benchmarks
-   **Well Tested**: Comprehensive test suite with extensive coverage

## Installation

```bash
go get github.com/dredfort42/go_config_reader
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    config "github.com/dredfort42/go_config_reader"
)

func main() {
    // Create a new config instance
    cfg, err := config.New()
    if err != nil {
        log.Fatal(err)
    }

    // Load configuration from file
    err = cfg.LoadFromFile("config.json", nil)
    if err != nil {
        log.Fatal(err)
    }

    // Get values with defaults
    port := cfg.GetInt("server_port", 8080)
    debug := cfg.GetBool("debug_mode", false)
    name := cfg.GetString("app_name", "MyApp")

    fmt.Printf("Server: %s running on port %d (debug: %v)\n", name, port, debug)
}
```

## Configuration Formats

### JSON Configuration

```json
{
    "server_port": 8080,
    "debug_mode": true,
    "database": {
        "host": "localhost",
        "port": 5432,
        "name": "myapp"
    },
    "features": ["auth", "api", "web"],
    "timeout": "30s"
}
```

### YAML Configuration

```yaml
server_port: 8080
debug_mode: true
database:
    host: localhost
    port: 5432
    name: myapp
features:
    - auth
    - api
    - web
timeout: 30s
```

### INI Configuration

```ini
# Server settings
server_port=8080
debug_mode=true
timeout=30

# Database settings
db_host=localhost
db_port=5432
db_name=myapp
```

## Data Types

The library supports automatic type conversion for:

-   **String**: `GetString(key, default...)`
-   **Integer**: `GetInt(key, default...)`
-   **Float64**: `GetFloat64(key, default...)`
-   **Boolean**: `GetBool(key, default...)`
-   **Duration**: `GetDuration(key, default...)` (supports "30s", "5m", "1h" format)
-   **String Slice**: `GetStringSlice(key, default...)` (supports arrays and comma-separated strings)

### Duration Examples

```go
// From config file
"timeout": "1h30m45s"
"retry_delay": "500ms"
"cache_ttl": "24h"

// From code
timeout := cfg.GetDuration("timeout", 30*time.Second)
retryDelay := cfg.GetDuration("retry_delay", 100*time.Millisecond)
```

## Environment Variable Override

Environment variables automatically override configuration file values when the key names match:

```bash
export server_port=9090
export debug_mode=true
```

```go
// Will use server_port=9090 from environment even if config file has different value
port := cfg.GetInt("server_port", 8080)
```

To disable environment variable override:

```go
opts := &config.LoadOptions{
    IgnoreEnv: true, // Disable environment variable override
}
err = cfg.LoadFromFile("config.json", opts)
```

## Utility Functions

### Configuration Management

```go
// Check if key exists
if cfg.Has("optional_feature") {
    feature := cfg.GetString("optional_feature")
    // ...
}

// Get all keys
keys := cfg.Keys()
fmt.Printf("Available keys: %v\n", keys)

// Get all configuration data
all := cfg.GetAll()
for key, value := range all {
    fmt.Printf("%s: %v\n", key, value)
}

// Load from map
data := map[string]interface{}{
    "key1": "value1",
    "key2": 42,
}
cfg.LoadFromMap(data)

// Clear all configuration
cfg.Clear()

// Check if configuration is empty
if cfg.IsEmpty() {
    fmt.Println("Configuration is empty")
}

// Get configuration size
size := cfg.Size()
fmt.Printf("Configuration has %d keys\n", size)

// Get string representation
fmt.Println(cfg.String())
```

## Thread Safety

The library is fully thread-safe and can be used in concurrent applications:

```go
cfg, _ := config.New()

// Safe to call from multiple goroutines
go func() {
    value := cfg.GetString("key")
    cfg.Set("runtime_key", "value")
}()

go func() {
    cfg.Set("another_key", "another_value")
    exists := cfg.Has("key")
}()
```

## API Reference

### Config Instance Methods

```go
// Creating a new config instance
cfg, err := config.New()

// Loading configuration
err = cfg.LoadFromFile(filePath, opts)
cfg.LoadFromMap(data)

// Getting values
cfg.GetString(key, defaultValue...)
cfg.GetInt(key, defaultValue...)
cfg.GetFloat64(key, defaultValue...)
cfg.GetBool(key, defaultValue...)
cfg.GetDuration(key, defaultValue...)
cfg.GetStringSlice(key, defaultValue...)

// Advanced getters
cfg.GetNestedMap(key)
cfg.GetNestedKeys(prefix)
cfg.GetAll()

// Setting and checking values
cfg.Set(key, value)
cfg.SetNestedDefaults(defaults)
cfg.Has(key)

// Utility methods
cfg.Keys()
cfg.Size()
cfg.IsEmpty()
cfg.Clear()
cfg.String()
```

### LoadOptions Structure

```go
type LoadOptions struct {
    Format         Format                     // Configuration file format (auto-detected if not specified)
    IgnoreEnv      bool                       // If true, skip environment variable override
    RequiredKeys   []string                   // Keys that must be present after loading
    DefaultValues  map[string]any             // Default values applied before loading file
    ValidationFunc func(map[string]any) error // Custom validation function
}
```

### Supported Formats

```go
const (
    FormatINI Format = iota
    FormatJSON
    FormatYAML
)
```

## Error Handling

```go
cfg, err := config.New()
if err != nil {
    log.Fatalf("Failed to create config: %v", err)
}

err = cfg.LoadFromFile("config.json", &config.LoadOptions{
    RequiredKeys: []string{"database_url"},
    ValidationFunc: func(data map[string]interface{}) error {
        // Custom validation logic
        return nil
    },
})
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}
```

## Examples

Comprehensive examples are available in the [`examples/`](examples/) directory:

-   **[Basic Examples](examples/main.go)**: Complete feature demonstration
-   **[HTTP Server Configuration](examples/server/main.go)**: Real-world server setup
-   **[Microservice Configuration](examples/microservice/main.go)**: Advanced microservice config

Each example includes detailed comments and demonstrates different aspects of the library.

## Requirements

-   Go 1.24 or later

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Author

**Dmitry Novikov**

-   Email: [dredfort.42@gmail.com](mailto:dredfort.42@gmail.com)
-   GitHub: [@dredfort42](https://github.com/dredfort42)
-   LinkedIn: [novikov-da](https://linkedin.com/in/novikov-da)
