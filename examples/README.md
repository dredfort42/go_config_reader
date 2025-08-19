# Go Configuration Reader - Examples

This directory contains comprehensive examples demonstrating all features and capabilities of the go_config_reader library.

## 📁 Examples Overview

### 1. Basic Examples (`main.go`)

**Comprehensive demonstration of all library features:**

-   Basic configuration loading
-   Different file formats (JSON, YAML, INI)
-   Type safety and getter methods
-   Nested configuration with dot notation
-   Default values
-   Environment variable override
-   Advanced loading options
-   Runtime configuration changes
-   Utility methods
-   Error handling
-   Thread safety
-   Real-world application configuration

**Run:** `go run main.go`

### 2. HTTP Server Configuration (`server/main.go`)

**Real-world HTTP server configuration example:**

-   Structured configuration with Go structs
-   Environment-specific defaults
-   Configuration validation
-   SSL/TLS configuration
-   Database connection settings
-   Server timeouts and limits
-   Non-sensitive configuration exposure

**Run:** `cd server && go run main.go`

### 3. Microservice Configuration (`microservice/main.go`)

**Advanced microservice configuration management:**

-   Environment-specific configurations (dev, staging, prod)
-   Service discovery configuration
-   Circuit breaker settings
-   Distributed tracing configuration
-   Health checks and metrics
-   Hot-reload simulation
-   Complex validation rules

**Run:** `cd microservice && go run main.go`

## 🚀 Quick Start

1. **Clone the repository:**

    ```bash
    git clone https://github.com/dredfort42/go_config_reader.git
    cd go_config_reader/examples
    ```

2. **Run the comprehensive examples:**

    ```bash
    go run main.go
    ```

3. **Try specific examples:**

    ```bash
    # HTTP Server Configuration
    cd server
    go run main.go

    # Microservice Configuration
    cd ../microservice
    go run main.go
    ```

## 📋 Configuration File Examples

The examples will create sample configuration files in various formats:

### JSON Configuration (`config.json`)

```json
{
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
}
```

### YAML Configuration (`config.yaml`)

```yaml
app_name: YamlApp
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
timeout: 45s
```

### INI Configuration (`config.ini`)

```ini
app_name=IniApp
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
reporting=false
```

## 🔧 Features Demonstrated

### Core Features

-   ✅ **Multiple Format Support**: JSON, YAML, INI
-   ✅ **Type Safety**: String, Int, Float64, Bool, Duration, StringSlice
-   ✅ **Nested Configuration**: Dot notation access (`database.host`)
-   ✅ **Default Values**: Fallback values for missing keys
-   ✅ **Environment Override**: Automatic environment variable priority
-   ✅ **Thread Safety**: Concurrent access protection
-   ✅ **Validation**: Custom validation functions
-   ✅ **Required Keys**: Mandatory configuration validation

### Advanced Features

-   ✅ **Runtime Changes**: Dynamic configuration updates
-   ✅ **Utility Methods**: Has, Keys, Size, IsEmpty, Clear, String
-   ✅ **Error Handling**: Comprehensive error scenarios
-   ✅ **Nested Defaults**: SetNestedDefaults for complex structures
-   ✅ **Configuration Parsing**: Integration with Go structs

### Real-World Scenarios

-   ✅ **Web Server Configuration**: Complete HTTP server setup
-   ✅ **Microservice Architecture**: Service discovery, tracing, circuits
-   ✅ **Environment Management**: Dev, staging, production configs
-   ✅ **Hot Reload**: Dynamic configuration updates
-   ✅ **Security**: Non-sensitive configuration exposure

## 🌍 Environment Variables

All examples support environment variable overrides. For example:

```bash
# Override configuration values
export server_port=9090
export debug_mode=true
export app_name="EnvironmentApp"

# Run examples with overrides
go run main.go
```

## 📊 Performance

The library is optimized for performance:

-   Concurrent reads using RWMutex
-   Minimal memory allocations
-   Fast type conversions
-   Efficient nested access

## 🛡️ Error Handling

Examples demonstrate proper error handling for:

-   File not found errors
-   Invalid JSON/YAML/INI syntax
-   Missing required keys
-   Validation failures
-   Type conversion errors
-   Nil configuration handling

## 📚 Learn More

-   **Main Documentation**: [../README.md](../README.md)
-   **API Reference**: See README.md for complete API documentation
-   **Contributing**: [../CONTRIBUTING.md](../CONTRIBUTING.md)
-   **License**: [../LICENSE](../LICENSE)

## 💡 Tips

1. **Start Simple**: Begin with basic configuration loading
2. **Use Defaults**: Always provide sensible default values
3. **Validate Early**: Use validation functions for critical settings
4. **Environment Variables**: Use for deployment-specific overrides
5. **Type Safety**: Use specific getter methods for type safety
6. **Nested Access**: Use dot notation for complex configurations
7. **Error Handling**: Always check for errors when loading configurations
