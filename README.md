# logdeb

[![Test](https://github.com/DarioEB/logdeb/actions/workflows/testing.yml/badge.svg)](https://github.com/DarioEB/logdeb/actions/workflows/testing.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/DarioEB/logdeb.svg)](https://pkg.go.dev/github.com/DarioEB/logdeb)

A simple multi-level logger for Go that writes to both console and JSON files. Each log level can be independently enabled and writes to its own file.

## Requirements

- Go 1.21+ (uses `log/slog` from standard library)

## Installation

```bash
go get github.com/DarioEB/logdeb
```

## Usage

### Basic usage with default config

```go
package main

import "github.com/DarioEB/logdeb"

func main() {
    logger, err := logdeb.New(logdeb.DefaultConfig())
    if err != nil {
        panic(err)
    }
    defer logger.Close()

    logger.Info("Application started")
    logger.Error("Something went wrong")
    logger.Warn("This is a warning")
    logger.Debug("Debug information")
}
```

### Custom configuration

```go
logger, err := logdeb.New(logdeb.Config{
    Dir:          "my_logs",    // Log directory (created if not exists)
    InfoEnabled:  true,         // Enable info.log
    ErrorEnabled: true,         // Enable error.log
    WarnEnabled:  false,        // Disable warn.log
    DebugEnabled: false,        // Disable debug.log
})
```

### Logging with attributes

```go
logger.Info("user login", "user_id", 123, "ip", "192.168.1.1")
logger.Error("database error", "query", "SELECT * FROM users", "error", err.Error())
```

Output in JSON file:
```json
{"time":"2026-01-25T10:30:00Z","level":"INFO","msg":"user login","user_id":123,"ip":"192.168.1.1"}
```

## Features

- Writes to both console (text) and files (JSON)
- Independent log levels: Info, Error, Warn, Debug
- Enable/disable each level individually
- Support for structured logging with key-value attributes
- Automatic directory creation
- Proper resource cleanup with `Close()`

## Log Files

When enabled, the following files are created in the configured directory:

| Level | File |
|-------|------|
| Info  | `info.log` |
| Error | `error.log` |
| Warn  | `warn.log` |
| Debug | `debug.log` |

## License

MIT License - see [LICENSE](LICENSE) file.
