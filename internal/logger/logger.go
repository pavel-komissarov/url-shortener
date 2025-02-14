package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(env string) *zap.Logger {
	var logger *zap.Logger

	var level zapcore.Level
	switch env {
	case "prod":
		level = zapcore.InfoLevel
	default:
		level = zapcore.DebugLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	core := zapcore.NewCore(encoder, os.Stdout, level)

	logger = zap.New(core)

	return logger
}
