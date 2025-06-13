# XGo - Go Utility Libraries

XGo provides a collection of utility packages for Go applications, including error handling, logging, type utilities,
and general-purpose utilities.

## Features

- **Custom Error Handling (xerror)**
    - Stack trace support
    - Error code management
    - Application name context
    - Timestamp tracking
    - gRPC integration

- **Advanced Logging (xlog)**
    - Structured logging with Zap
    - OpenTelemetry integration
    - Pretty printing for development
    - Context-aware logging
    - Support for Gin web framework

- **Type Utilities (xtype)**
    - Date handling with custom JSON marshaling
    - Phone number validation and formatting
    - Secure string encryption and hashing
    - Safe type conversions
    - Common type definitions

- **General Utilities (xutil)**
    - Generic mapping functions
    - Asynchronous processing
    - Enum handling
    - Type-safe conversions

## Installation

```bash
go get github.com/kurzgesagtz/xgo
```

## Usage

### Error Handling

```go
import "github.com/kurzgesagtz/xgo/xerror"

// Create a new error with a code
err := xerror.NewError(xerror.ErrCodeNotFound, 
    xerror.WithMessage("User not found"))

// Check error code
if xerror.IsErrorCode(err, xerror.ErrCodeNotFound) {
    // Handle not found error
}
```

### Logging

```go
import "github.com/kurzgesagtz/xgo/xlog"

// Simple logging
xlog.Info().Msg("Application started")

// Structured logging
xlog.Debug().
    Field("user_id", 123).
    Field("action", "login").
    Msg("User logged in")

// Error logging
err := someFunction()
if err != nil {
    xlog.Error().Err(err).Msg("Operation failed")
}
```
