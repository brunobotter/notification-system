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
}

type loggerZap struct {
	appName      string
	level        string
	logger       *zap.Logger
	commonFields []any
}

func NewLoggerZap(appName string) Logger {
	var config zap.Config
	var zapLogger *zap.Logger

	config.Encoding = "json"
	config.EncoderConfig = buildEncondingConfig()

	config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	zapLogger, _ = config.Build()

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

func buildEncondingConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "key",
		TimeKey:        "timestamp",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (l *loggerZap) InfoF(format string, args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Infof(format, args...)
}

func (l *loggerZap) Info(args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Info(args...)
}

func (l *loggerZap) ErrorF(format string, args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Errorf(format, args...)
}

func (l *loggerZap) Error(format string, args ...interface{}) {
	defer l.logger.Sync()
	l.logger.Sugar().With(l.commonFields...).Error(args...)
}

func (l *loggerZap) Log(msg string) {
	l.Info(msg)
}
