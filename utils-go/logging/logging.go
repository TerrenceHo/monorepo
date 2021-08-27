package logging

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerType string

const (
	DevLogger  LoggerType = "dev-logger"
	ProdLogger LoggerType = "prod-logger"
	TestLogger LoggerType = "test-logger"
)

var globalLoggerLock sync.RWMutex
var globalLogger Logger

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

// SetGlobalLogger sets the global logger
func SetGlobalLogger(logger Logger) {
	globalLoggerLock.Lock()
	globalLogger = logger
	globalLoggerLock.Unlock()
}

// GetGlobalLogger returns the global logger
func GetGlobalLogger() Logger {
	globalLoggerLock.RLock()
	logger := globalLogger
	globalLoggerLock.RUnlock()
	return logger
}

// ConfigureLogger takes in a logger type, and returns the logger for that
// environment. For tests, it returns a Nop logger. For development, it returns
// a logger that writes to the console in a human readable way, with caller and
// stack traces enabled for warn level errors. For production, it removes the
// callers and stack traces (except in Error level logs) in favor of
// speed/performance.
func ConfigureLogger(loggerType LoggerType) (Logger, error) {
	var logger Logger
	var err error
	switch loggerType {
	case TestLogger:
		logger = zap.NewNop()
	case DevLogger:
		logger, err = zap.NewDevelopment(
			zap.ErrorOutput(os.Stderr),
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	case ProdLogger:
		logger, err = zap.NewProduction(
			zap.ErrorOutput(os.Stderr),
		)
	}
	return logger, err
}

func Debug(msg string, fields ...zap.Field) {
	GetGlobalLogger().Debug(msg, fields...)
}
func Info(msg string, fields ...zap.Field) {
	GetGlobalLogger().Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	GetGlobalLogger().Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	GetGlobalLogger().Error(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}
