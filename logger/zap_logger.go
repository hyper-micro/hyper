package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	zap *zap.SugaredLogger
}

var zapLevel = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
}

const (
	EncoderJSON    = "json"
	EncoderConsole = "console"
)

type ZapLoggerConfig struct {
	Level   string
	Writer  []io.Writer
	Encoder string
	Caller  bool
	Fn      bool
}

func NewZapLogger(conf ZapLoggerConfig) Logger {
	lvl := zap.InfoLevel
	level := ParseLevel(conf.Level)
	if l, ok := zapLevel[level]; ok {
		lvl = l
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lvl)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if conf.Fn {
		encoderConfig.FunctionKey = "fn"
	}
	encoderConfig.MessageKey = "msg"

	var encoder zapcore.Encoder
	switch conf.Encoder {
	case EncoderConsole:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	var syncer []zapcore.WriteSyncer
	for _, w := range conf.Writer {
		syncer = append(syncer, zapcore.AddSync(w))
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(syncer...),
		atomicLevel,
	)
	z := zap.New(
		core,
		zap.WithCaller(conf.Caller),
		zap.AddCallerSkip(1),
	).Sugar()

	return &zapLogger{
		zap: z,
	}
}

func (l *zapLogger) z() *zap.SugaredLogger {
	if l.zap == nil {
		panic("logger: Zap.zap not initialized")
	}
	return l.zap
}

func (l *zapLogger) Debug(args ...interface{}) {
	l.z().Debug(args...)
}

func (l *zapLogger) Info(args ...interface{}) {
	l.z().Info(args...)
}

func (l *zapLogger) Warn(args ...interface{}) {
	l.z().Warn(args...)
}

func (l *zapLogger) Error(args ...interface{}) {
	l.z().Error(args...)
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.z().Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.z().Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.z().Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.z().Errorf(format, args...)
}
