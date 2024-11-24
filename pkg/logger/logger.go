package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	if strings.ToLower(os.Getenv("APP_ENV")) == "development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder

	config := zap.Config{
		// Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Level:             zap.NewAtomicLevelAt(getLoggerLevel()),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}

func getLoggerLevel() zapcore.Level {
	logLevel := strings.ToLower(os.Getenv("APP_LOG_LEVEL"))
	var level zapcore.Level
	switch {
	case logLevel == "debug":
		level = zapcore.DebugLevel
	case logLevel == "warn":
		level = zapcore.WarnLevel
	case logLevel == "error":
		level = zapcore.ErrorLevel
	case logLevel == "dpanic":
		level = zapcore.DPanicLevel
	case logLevel == "panic":
		level = zapcore.PanicLevel
	case logLevel == "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.InfoLevel
	}
	return level
}
