package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	SetCommonField(commonFields map[string]any)
	InfoF(format string, args ...interface{})
	Info(args ...interface{})
	ErrorF(format string, args ...interface{})
	Error(format string, args ...interface{})
	Log(msg string)
	Sync()
}

type loggerZap struct {
	appName      string
	level        string
	logger       *zap.Logger
	commonFields []any
}

func NewLoggerZap(appName string) Logger {
	config := zap.NewProductionConfig()
	config.Encoding = "json"
	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	zapLogger, _ := config.Build()

	j := &loggerZap{
		appName:      appName,
		level:        config.Level.String(),
		logger:       zapLogger,
		commonFields: []any{},
	}

	j.SetCommonField(map[string]any{
		"application_name": appName,
	})

	return j
}

func (l *loggerZap) SetCommonField(commanFields map[string]any) {
	for key, value := range commanFields {
		l.commonFields = append(l.commonFields, zap.Any(key, value))
	}
}

func (l *loggerZap) InfoF(format string, args ...interface{}) {
	l.logger.Sugar().With(l.commonFields...).Infof(format, args...)
}

func (l *loggerZap) Info(args ...interface{}) {
	l.logger.Sugar().With(l.commonFields...).Info(args...)
}

func (l *loggerZap) ErrorF(format string, args ...interface{}) {
	l.logger.Sugar().With(l.commonFields...).Errorf(format, args...)
}

func (l *loggerZap) Error(format string, args ...interface{}) {
	l.logger.Sugar().With(l.commonFields...).Error(args...)
}

func (l *loggerZap) Log(msg string) {
	l.Info(msg)
}

func (l *loggerZap) Sync() {
	_ = l.logger.Sync()
}
