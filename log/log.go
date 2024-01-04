package log

import (
	"io"
	"strings"

	"github.com/hyper-micro/hyper/log/writer"
	"github.com/hyper-micro/hyper/slice"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type Level int8

const (
	NoneLevel Level = iota
	DebugLevel
	InfoLevel
	WarnLevel
	ErrorLevel
)

var unmarshalLevelText = map[string]Level{
	"debug": DebugLevel,
	"info":  InfoLevel,
	"warn":  WarnLevel,
	"error": ErrorLevel,
}

func ParseLevel(text string) Level {
	text = strings.ToLower(text)
	lvl, _ := unmarshalLevelText[text]
	return lvl
}

type Config struct {
	Output         []string
	FilePath       string
	Level          string
	MaxRotatedSize int
	MaxRetainDay   int
	MaxRetainFiles int
	LocalTime      bool
	Encoder        string
	Caller         bool
	Fn             bool
}

func NewLogger(conf Config) Logger {
	if len(conf.Output) == 0 {
		conf.Output = append(conf.Output, "stdout")
	}

	var writers []io.Writer

	if slice.ContainsString("file", conf.Output) {
		if conf.FilePath != "" {
			fileWriter := writer.NewLumberJackWriter(writer.LumberJackConfig{
				FilePath:       conf.FilePath,
				MaxRotatedSize: conf.MaxRotatedSize,
				MaxRetainDay:   conf.MaxRetainDay,
				MaxRetainFiles: conf.MaxRetainFiles,
				LocalTime:      conf.LocalTime,
			})
			writers = append(writers, fileWriter)
		}
	}
	if slice.ContainsString("stdout", conf.Output) {
		writers = append(writers, writer.NewStdoutWriter())
	}

	driver := NewZapLogger(ZapLoggerConfig{
		Level:   conf.Level,
		Writer:  writers,
		Encoder: conf.Encoder,
		Caller:  conf.Caller,
		Fn:      conf.Fn,
	})
	return driver
}

var logger = NewLogger(Config{
	Level:   "debug",
	Encoder: "console",
})

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}
