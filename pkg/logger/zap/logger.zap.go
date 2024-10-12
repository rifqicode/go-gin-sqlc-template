package zap

import (
	zaplogger "go.uber.org/zap"
	zapcore "go.uber.org/zap/zapcore"
)

var logger *zaplogger.Logger

func init() {
	logger, _ = zaplogger.NewProduction()
}

func Info(message string, fields ...interface{}) {
	logger.Info(message, convertToZapFields(fields)...)
}

func Debug(message string, fields ...interface{}) {
	logger.Debug(message, convertToZapFields(fields)...)
}

func Warn(message string, fields ...interface{}) {
	logger.Warn(message, convertToZapFields(fields)...)
}

func Error(message string, fields ...interface{}) {
	logger.Error(message, convertToZapFields(fields)...)
}

func Fatal(message string, fields ...interface{}) {
	logger.Fatal(message, convertToZapFields(fields)...)
}

func Panic(message string, fields ...interface{}) {
	logger.Panic(message, convertToZapFields(fields)...)
}

// Sugar logger methods for convenience

func Infof(template string, args ...interface{}) {
	logger.Sugar().Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Sugar().Debugf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Sugar().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Sugar().Fatalf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	logger.Sugar().Panicf(template, args...)
}

// Helper function to convert interface{} to zapcore.Field
func convertToZapFields(fields []interface{}) []zapcore.Field {
	zapFields := make([]zapcore.Field, 0, len(fields))
	for _, field := range fields {
		switch f := field.(type) {
		case zapcore.Field:
			zapFields = append(zapFields, f)
		case map[string]interface{}:
			for k, v := range f {
				zapFields = append(zapFields, zaplogger.Any(k, v))
			}
		default:
			zapFields = append(zapFields, zaplogger.Any("value", f))
		}
	}
	return zapFields
}
