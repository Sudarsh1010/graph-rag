package di

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(config *Config) (*zap.Logger, error) {
	// Parse log level
	var zapLevel zapcore.Level
	switch config.LogLevel {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}

	// Build zap config based on format
	var cfgZap zap.Config
	if config.LogFormat == "json" {
		cfgZap = zap.NewProductionConfig()
	} else {
		cfgZap = zap.NewDevelopmentConfig()
	}

	cfgZap.Level = zap.NewAtomicLevelAt(zapLevel)

	return cfgZap.Build()
}
