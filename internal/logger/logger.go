package logger

import (
	"errors"
	"fmt"
	"github.com/neatplex/nightel-core/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"syscall"
)

type Logger struct {
	e     *zap.Logger
	close func()
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.e.Debug(message, fields...)
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	l.e.Info(message, fields...)
}

func (l *Logger) Warn(message string, fields ...zap.Field) {
	l.e.Warn(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	l.e.Error(message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.close()
	l.e.Error(message, fields...)
}

func (l *Logger) With(fields ...zap.Field) *zap.Logger {
	return l.e.With(fields...)
}

func (l *Logger) Close() {
	if err := l.e.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		l.e.Warn("cannot close the log", zap.Error(err))
	} else {
		l.e.Debug("log closed successfully")
	}
}

func New(c *config.Config, close func()) (logger *Logger, err error) {
	level := zap.NewAtomicLevel()
	if err = level.UnmarshalText([]byte(c.Logger.Level)); err != nil {
		return nil, errors.New(fmt.Sprintf("invalid log level %s, err: %v", c.Logger.Level, err))
	}

	engine, err := zap.Config{
		Level:             level,
		Development:       false,
		Encoding:          "json",
		DisableStacktrace: true,
		DisableCaller:     true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
			EncodeDuration: zapcore.StringDurationEncoder,
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			NameKey:        "key",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
		},
	}.Build()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("cannot build logger, err: %v", err))
	}

	return &Logger{e: engine, close: close}, nil
}
