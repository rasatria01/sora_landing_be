package observability

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type LogFields map[string]interface{}

// Logger wraps zap logger with simplified methods
type Logger struct {
	log *zap.Logger
}

// NewLogger creates a new logger instance
func NewLogger() (*Logger, error) {
	log, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &Logger{log: log}, nil
}

// WithContext adds context values to log fields
func (l *Logger) WithContext(ctx context.Context) *Logger {
	// Extract trace ID or request ID if present
	if requestID := ctx.Value("request_id"); requestID != nil {
		l.log = l.log.With(zap.String("request_id", requestID.(string)))
	}
	return l
}

// Info logs info level message with fields
func (l *Logger) Info(msg string, fields LogFields) {
	zapFields := toZapFields(fields)
	l.log.Info(msg, zapFields...)
}

// Error logs error level message with fields
func (l *Logger) Error(msg string, err error, fields LogFields) {
	zapFields := toZapFields(fields)
	zapFields = append(zapFields, zap.Error(err))
	l.log.Error(msg, zapFields...)
}

// toZapFields converts LogFields to zap.Field slice
func toZapFields(fields LogFields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		switch val := v.(type) {
		case string:
			zapFields = append(zapFields, zap.String(k, val))
		case int:
			zapFields = append(zapFields, zap.Int(k, val))
		case bool:
			zapFields = append(zapFields, zap.Bool(k, val))
		case time.Duration:
			zapFields = append(zapFields, zap.Duration(k, val))
		default:
			zapFields = append(zapFields, zap.Any(k, v))
		}
	}
	return zapFields
}
