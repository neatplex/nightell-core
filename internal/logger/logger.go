package logger

import (
	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"syscall"
)

type Logger struct {
	level       string
	format      string
	development bool
	Engine      *zap.Logger
}

func (l *Logger) Init() (err error) {
	level := zap.NewAtomicLevel()
	if err = level.UnmarshalText([]byte(l.level)); err != nil {
		return errors.Errorf("logger: invalid level %s, err: %v", l.level, err)
	}

	l.Engine, err = zap.Config{
		Level:             level,
		Development:       l.development,
		Encoding:          "json",
		DisableStacktrace: false,
		DisableCaller:     false,
		OutputPaths:       []string{"./storage/logs/output.log"},
		ErrorOutputPaths:  []string{"./storage/logs/output.log"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			EncodeTime:     zapcore.TimeEncoderOfLayout(l.format),
			EncodeDuration: zapcore.StringDurationEncoder,
			LevelKey:       "level",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			NameKey:        "key",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "message",
			LineEnding:     zapcore.DefaultLineEnding,
		},
	}.Build()

	return errors.WithStack(err)
}

func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Engine.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Engine.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Engine.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Engine.Error(msg, fields...)
}

func (l *Logger) With(fields ...zap.Field) *zap.Logger {
	return l.Engine.With(fields...)
}

func (l *Logger) Close() {
	if err := l.Engine.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		l.Engine.Error("logger: failed to close", zap.Error(err))
	} else {
		l.Engine.Info("logger: closed successfully")
	}
}

func New(level, format string, development bool) (logger *Logger) {
	return &Logger{Engine: nil, level: level, format: format, development: development}
}
