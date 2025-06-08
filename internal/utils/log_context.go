package utils

import (
	"go.uber.org/zap"
)

// LogContext holds contextual fields for structured logging
type LogContext struct {
	logger *zap.SugaredLogger
}

// NewLogContext creates a new logging context with fields
func NewLogContext(fields ...interface{}) *LogContext {
	if Logger == nil {
		initLogger()
	}
	return &LogContext{
		logger: SugaredLogger.With(fields...),
	}
}

// Debug logs a debug message with context
func (lc *LogContext) Debug(format string, v ...interface{}) {
	lc.logger.Debugf(format, v...)
}

// Info logs an info message with context
func (lc *LogContext) Info(format string, v ...interface{}) {
	lc.logger.Infof(format, v...)
}

// Warn logs a warning message with context
func (lc *LogContext) Warn(format string, v ...interface{}) {
	lc.logger.Warnf(format, v...)
}

// Error logs an error message with context
func (lc *LogContext) Error(format string, v ...interface{}) {
	lc.logger.Errorf(format, v...)
}

// Fatal logs a fatal message with context and exits
func (lc *LogContext) Fatal(format string, v ...interface{}) {
	lc.logger.Fatalf(format, v...)
}

// With adds additional context fields
func (lc *LogContext) With(fields ...interface{}) *LogContext {
	return &LogContext{
		logger: lc.logger.With(fields...),
	}
}
