package logger

import (
	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/logger"
)

type Provider interface {
	Into() logger.Logger
}

type loggerProvider struct {
	logger logger.Logger
}

func NewProvider(conf config.Config) (Provider, func(), error) {
	instance := logger.NewLogger(logger.Config{
		Output:         []string{"file"},
		FilePath:       conf.GetStringOrDefault("log.logger.path", "logs"),
		Level:          conf.GetStringOrDefault("log.logger.level", "error"),
		MaxRotatedSize: conf.GetInt("log.logger.rotatedSize"),
		MaxRetainDay:   conf.GetInt("log.logger.retainDay"),
		MaxRetainFiles: conf.GetInt("log.logger.retainFiles"),
		Encoder:        "json",
		Caller:         true,
	})
	return &loggerProvider{logger: instance}, func() {}, nil
}

func (p *loggerProvider) Into() logger.Logger {
	return p.logger
}
