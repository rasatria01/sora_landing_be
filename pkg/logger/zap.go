package logger

import (
	"log"
	"sora_landing_be/pkg/config"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once sync.Once
	Log  *ZapLogger
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(config config.Logger) {
	var logger *zap.Logger
	once.Do(func() {
		level, err := zapcore.ParseLevel(config.LogLevel)
		if err != nil {
			log.Fatal("failed to parse log level: %w", err)
		}

		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(level)
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		zapConfig.Encoding = config.Encoding
		zapConfig.DisableStacktrace = true

		logger, err = zapConfig.Build()
		logger = logger.WithOptions(zap.AddCallerSkip(1))
		if err != nil {
			log.Fatal(err)
		}
		Log = &ZapLogger{
			logger: logger,
		}
	})
}

// Info logs an informational message.
func (z *ZapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

// Error logs an error message.
func (z *ZapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

// Debug logs a debug message.
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

// Sync flushes any buffered log entries.
func (z *ZapLogger) Sync() error {
	return z.logger.Sync()
}

func (z *ZapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

// Fatal logs a fatal message and then calls os.Exit(1)
func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}
