# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
