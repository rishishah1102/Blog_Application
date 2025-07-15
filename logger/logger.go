package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

// NewLogger initializes a custom zap logger
func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	globalLogger = logger
	return logger, nil
}

// WrapError logs and wraps an existing error
func WrapError(err error, msg string, fields ...zap.Field) error {
	fields = append(fields, zap.Error(err))
	globalLogger.Error(msg, fields...)
	return fmt.Errorf("%s: %w", msg, err)
}
