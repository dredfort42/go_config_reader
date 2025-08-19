# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-08-19

### Added

-   Comprehensive examples package with detailed demonstrations
-   HTTP server configuration example (`examples/server/main.go`)
-   Microservice configuration example (`examples/microservice/main.go`)
-   Advanced getter methods: `GetNestedMap()` and `GetNestedKeys()`
-   `SetNestedDefaults()` method for setting default values for nested keys
-   Examples documentation in `examples/README.md`

### Fixed

-   Corrected README to accurately reflect current API (removed non-existent global functions)
-   Updated documentation to match actual implementation
-   Improved examples with better error handling
-   Fixed all code examples to work with current API

## [1.0.0] - 2025-08-07

### Added

-   Initial release of go-config-reader
-   Support for JSON, YAML, and INI configuration formats
-   Thread-safe configuration management with RWMutex
-   Type-safe getters with automatic type conversion
-   Environment variable override support
-   Default values support
-   Configuration validation with custom validation functions
-   Required key checking
-   Global API for backward compatibility
-   Comprehensive test suite with benchmarks
-   Example configurations and usage demonstrations

### Features

-   **Multiple Format Support**: JSON, YAML, and INI files
-   **Environment Variable Override**: Automatic environment variable priority
-   **Thread-Safe**: Concurrent access protection with RWMutex
-   **Type Safety**: Strong typing with automatic type conversion
-   **Default Values**: Built-in support for default values
-   **Validation**: Custom validation functions and required key checking
-   **Backward Compatibility**: Maintains compatibility with existing code
-   **Zero Dependencies**: Minimal external dependencies (only gopkg.in/yaml.v3)

### Dependencies

-   Go 1.24+
-   gopkg.in/yaml.v3 v3.0.1
-   github.com/stretchr/testify v1.10.0 (test only)
