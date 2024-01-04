package writer

import (
	"io"

	"github.com/natefinch/lumberjack"
)

type LumberJackConfig struct {
	FilePath       string
	MaxRotatedSize int
	MaxRetainDay   int
	MaxRetainFiles int
	LocalTime      bool
}

func NewLumberJackWriter(conf LumberJackConfig) io.Writer {
	return &lumberjack.Logger{
		Filename:   conf.FilePath,
		MaxSize:    conf.MaxRotatedSize,
		MaxAge:     conf.MaxRetainDay,
		MaxBackups: conf.MaxRetainFiles,
		LocalTime:  conf.LocalTime,
	}
}
