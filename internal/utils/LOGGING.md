# Logging with Zap

This project uses Uber's Zap logger for high-performance, structured logging.

## Configuration

Logging can be configured in the `config.yaml` file:

```yaml
logger:
  level: "info"           # debug, info, warn, error, fatal
  encoding: "json"        # json or console
  output_paths:
    - "stdout"
    - "./logs/app.log"
  error_output_paths:
    - "stderr"
    - "./logs/error.log"
```

## Basic Usage

```go
import "github.com/dksensei/letsnormalizeit/internal/utils"

// Simple logging
utils.Debug("Debug message: %s", "details")
utils.Info("Info message: %s", "details")
utils.Warn("Warning message: %s", "details")
utils.Error("Error message: %s", "details")
utils.Fatal("Fatal message: %s", "details") // Will exit the application

// With structured fields
logger := utils.With("userID", 123, "requestID", "abc-123")
logger.Info("Processing request")
```

## Using LogContext

The `LogContext` struct provides a way to carry structured logging fields:

```go
// Create a logger with context
ctx := utils.NewLogContext("userID", 123, "requestID", "abc-123")
ctx.Info("Processing request")

// Add more context
enrichedCtx := ctx.With("action", "login")
enrichedCtx.Info("User logged in")
```

## Performance Considerations

Zap is designed to be extremely fast and efficient, with minimal memory allocations.
For high-performance logging needs, prefer structured logging over string formatting.
