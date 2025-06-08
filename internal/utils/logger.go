package utils

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel represents the logging level
type LogLevel int

const (
	// DEBUG level
	DEBUG LogLevel = iota
	// INFO level
	INFO
	// WARN level
	WARN
	// ERROR level
	ERROR
	// FATAL level
	FATAL
)

var (
	// Level sets the current log level
	Level = INFO

	// Logger is the zap logger instance
	Logger *zap.Logger

	// SugaredLogger provides a more convenient API than Logger
	SugaredLogger *zap.SugaredLogger

	once sync.Once
)

// LevelNames maps log levels to their string representations
var LevelNames = map[LogLevel]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
	FATAL: "FATAL",
}

// SetLogLevel sets the logging level
func SetLogLevel(level LogLevel) {
	Level = level
	initLogger() // Re-initialize logger with new level
}

// getZapLevel converts our LogLevel to zapcore.Level
func getZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	case WARN:
		return zapcore.WarnLevel
	case ERROR:
		return zapcore.ErrorLevel
	case FATAL:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// getZapLevelFromString converts a string level to zapcore.Level
func getZapLevelFromString(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// initLogger initializes the zap logger
func initLogger() {
	once.Do(func() {
		// Create encoder configuration
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		// Create core
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			getZapLevel(Level),
		)

		// Create logger
		Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
		SugaredLogger = Logger.Sugar()
	})
}

// InitWithConfig initializes the logger with custom configuration
func InitWithConfig(config map[string]interface{}) {
	level := config["level"].(string)
	encoding := config["encoding"].(string)
	outputPaths := config["output_paths"].([]string)
	errorOutputPaths := config["error_output_paths"].([]string)

	zapLevel := getZapLevelFromString(level)

	// Convert our LogLevel based on the zapLevel
	switch zapLevel {
	case zapcore.DebugLevel:
		Level = DEBUG
	case zapcore.InfoLevel:
		Level = INFO
	case zapcore.WarnLevel:
		Level = WARN
	case zapcore.ErrorLevel:
		Level = ERROR
	case zapcore.FatalLevel:
		Level = FATAL
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Create encoder based on config
	// var encoder zapcore.Encoder
	// if encoding == "json" {
	// 	encoder = zapcore.NewJSONEncoder(encoderConfig)
	// } else {
	// 	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	// }

	// Set up output paths
	// This ensures that log files are created with the right permissions
	for _, path := range outputPaths {
		if path != "stdout" && path != "stderr" {
			dir := filepath.Dir(path)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
		}
	}

	for _, path := range errorOutputPaths {
		if path != "stdout" && path != "stderr" {
			dir := filepath.Dir(path)
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}
		}
	}

	// Create logger config
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapLevel),
		Encoding:         encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
	}

	// Build logger
	var err error
	Logger, err = cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	SugaredLogger = Logger.Sugar()
}

// Init initializes the logger
func Init(level LogLevel) {
	Level = level
	initLogger()
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	if Logger == nil {
		initLogger()
	}
	SugaredLogger.Debugf(format, v...)
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	if Logger == nil {
		initLogger()
	}
	SugaredLogger.Infof(format, v...)
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	if Logger == nil {
		initLogger()
	}
	SugaredLogger.Warnf(format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	if Logger == nil {
		initLogger()
	}
	SugaredLogger.Errorf(format, v...)
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	if Logger == nil {
		initLogger()
	}
	SugaredLogger.Fatalf(format, v...)
}

// With returns a Logger with the specified fields added to the log context
func With(fields ...interface{}) *zap.SugaredLogger {
	if Logger == nil {
		initLogger()
	}
	return SugaredLogger.With(fields...)
}

// Sync flushes any buffered log entries
func Sync() error {
	if Logger == nil {
		return nil
	}
	return Logger.Sync()
}

// Log logs a message with the specified level (for backward compatibility)
func Log(level LogLevel, format string, v ...interface{}) {
	if level < Level {
		return
	}

	if Logger == nil {
		initLogger()
	}

	switch level {
	case DEBUG:
		SugaredLogger.Debugf(format, v...)
	case INFO:
		SugaredLogger.Infof(format, v...)
	case WARN:
		SugaredLogger.Warnf(format, v...)
	case ERROR:
		SugaredLogger.Errorf(format, v...)
	case FATAL:
		SugaredLogger.Fatalf(format, v...)
	}
}
