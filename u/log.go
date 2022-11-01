package u

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Log interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

var Logger *SwitchboardLogger

type SwitchboardLogger struct {
	ZapLogger *zap.Logger
	level     zapcore.Level
}

func InitializeLogger(environment string) error {
	logLevel, err := getLogLevelFromEnv(environment)
	if err != nil {
		return err
	}
	Logger = New(os.Stdout, logLevel)
	return nil
}

func New(writer io.Writer, level zapcore.Level) *SwitchboardLogger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()

	var encoder zapcore.Encoder
	if level == zap.DebugLevel {
		encoder = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		level,
	)
	logger := &SwitchboardLogger{
		ZapLogger: zap.New(core),
		level:     level,
	}
	return logger
}

func (l *SwitchboardLogger) Debug(msg string, fields ...zap.Field) {
	coloredText, _ := ColorTextByName(msg, "blue")
	l.ZapLogger.Debug(coloredText, fields...)
}

func (l *SwitchboardLogger) Info(msg string, fields ...zap.Field) {
	l.ZapLogger.Info(msg, fields...)
}

func (l *SwitchboardLogger) Warn(msg string, fields ...zap.Field) {
	coloredText, _ := ColorTextByName(msg, "yellow")
	l.ZapLogger.Warn(coloredText, fields...)
}
func (l *SwitchboardLogger) Error(msg string, fields ...zap.Field) {
	l.ZapLogger.Error(msg, fields...)
}
func (l *SwitchboardLogger) DPanic(msg string, fields ...zap.Field) {
	l.ZapLogger.DPanic(msg, fields...)
}
func (l *SwitchboardLogger) Panic(msg string, fields ...zap.Field) {
	l.ZapLogger.Panic(msg, fields...)
}

func (l *SwitchboardLogger) Fatal(msg string, fields ...zap.Field) {
	l.ZapLogger.Fatal(msg, fields...)
}

func (l *SwitchboardLogger) Sync() error {
	return l.ZapLogger.Sync()
}

func getLogLevelFromEnv(environment string) (zapcore.Level, error) {
	switch environment {
	case "development":
		return zap.DebugLevel, nil
	case "production":
		return zap.InfoLevel, nil
	case "testing":
		return zap.DebugLevel, nil
	default:
		return -1, fmt.Errorf("invalid environment, cannot deduce log level")
	}
}
