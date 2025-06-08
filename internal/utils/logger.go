package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
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

	// Logger is the default logger
	Logger = log.New(os.Stdout, "", 0)
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
}

// Log logs a message with the specified level
func Log(level LogLevel, format string, v ...interface{}) {
	if level < Level {
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")

	// Get file and line information
	_, file, line, ok := runtime.Caller(1)
	fileInfo := ""
	if ok {
		fileInfo = fmt.Sprintf(" [%s:%d]", file, line)
	}

	levelName := LevelNames[level]
	message := fmt.Sprintf(format, v...)

	Logger.Printf("%s %s%s: %s", now, levelName, fileInfo, message)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug logs a debug message
func Debug(format string, v ...interface{}) {
	Log(DEBUG, format, v...)
}

// Info logs an info message
func Info(format string, v ...interface{}) {
	Log(INFO, format, v...)
}

// Warn logs a warning message
func Warn(format string, v ...interface{}) {
	Log(WARN, format, v...)
}

// Error logs an error message
func Error(format string, v ...interface{}) {
	Log(ERROR, format, v...)
}

// Fatal logs a fatal message and exits
func Fatal(format string, v ...interface{}) {
	Log(FATAL, format, v...)
}
