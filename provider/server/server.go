package server

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/logger"
)

type Provider interface {
	Run() error
	RegServe(f RegServeHandler) error
}

type App interface {
	Run() error
	Shutdown() error
}

type CleanUpHandler func()

type RegServeOption struct {
	Conf config.Config
}

type RegServeHandler func(RegServeOption) (App, CleanUpHandler, error)

type serverProvider struct {
	opt             Option
	apps            []App
	cleanUps        []func()
	flagSet         *flag.FlagSet
	configFileFlag  string
	showHelpFlag    bool
	showVersionFlag bool
	inShutdown      bool
	regServeOption  RegServeOption
}

type Option struct {
	AppName      string
	Version      string
	BuildCommit  string
	BuildDate    string
	ShutdownSigs []os.Signal
}

func NewProvider(opt Option) (Provider, func(), error) {
	srv := &serverProvider{
		opt:     opt,
		flagSet: flag.NewFlagSet(os.Args[0], flag.ContinueOnError),
	}

	if err := srv.init(); err != nil {
		return nil, nil, err
	}

	return srv, nil, nil
}

func (s *serverProvider) RegServe(f RegServeHandler) error {
	app, cleanUp, err := f(s.regServeOption)
	if err != nil {
		return err
	}

	s.apps = append(s.apps, app)
	s.cleanUps = append(s.cleanUps, cleanUp)
	return nil
}

func (s *serverProvider) Run() error {

	if len(s.opt.ShutdownSigs) > 0 {
		shutdownSignChan := make(chan os.Signal, 1)
		signal.Notify(shutdownSignChan, s.opt.ShutdownSigs...)
		go func() {
			recSign := <-shutdownSignChan
			s.stdLoggerPrint("Receive signal: %v", recSign)
			s.shutdown()
		}()
	}

	defer func() {
		s.shutdown()
		s.stdLoggerPrint("Server stopped, Bye!")
	}()

	s.stdLoggerPrint("Load config file: %s", s.configFileFlag)
	s.stdLoggerPrint("Version: %s, %s, %s", s.opt.Version, s.opt.BuildCommit, s.opt.BuildDate)
	s.stdLoggerPrint("Pid: %v", os.Getpid())
	s.stdLoggerPrint("Signal.Notify: %v", s.opt.ShutdownSigs)

	var appErr = make(chan error)

	for _, app := range s.apps {
		go func(app App) {
			appErr <- app.Run()
		}(app)
	}

	return <-appErr
}

func (s *serverProvider) shutdown() {
	if s.inShutdown {
		return
	}

	s.inShutdown = true

	for _, app := range s.apps {
		s.stdLoggerPrint("shutting down")
		if err := app.Shutdown(); err != nil {
			s.stdErrLoggerPrint("shutdown failed: %v", err)
		}
	}
	for _, cleanUp := range s.cleanUps {
		cleanUp()
	}
}

func (s *serverProvider) init() error {
	s.flagSet.Usage = func() {}
	s.flagSet.SetOutput(io.Discard)
	s.flagSet.StringVar(&s.configFileFlag, "c", "./conf/", "set configure file path")
	s.flagSet.StringVar(&s.configFileFlag, "config", "./conf/", "set configure file path")
	s.flagSet.BoolVar(&s.showHelpFlag, "h", false, "show help")
	s.flagSet.BoolVar(&s.showHelpFlag, "help", false, "show help")
	s.flagSet.BoolVar(&s.showVersionFlag, "v", false, "show version")
	s.flagSet.BoolVar(&s.showVersionFlag, "version", false, "show version")

	if err := s.flagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	if s.showHelpFlag {
		s.flagUsage()
		os.Exit(0)
	}

	if s.showVersionFlag {
		fmt.Printf("Version %s, build %s, %s\n", s.opt.Version, s.opt.BuildCommit, s.opt.BuildDate)
		os.Exit(0)
	}

	conf, err := config.New(config.PathTypePath, s.configFileFlag)
	if err != nil {
		return err
	}
	s.regServeOption.Conf = conf

	return nil
}

func (s *serverProvider) flagUsage() {
	fmt.Printf(`
USAGE:
   app [options]

A self-sufficient runtime for containers

OPTIONS:
   --config value, -c value  set configure file path (default: "./conf/config.yaml")
   --version, -v             show version (default: false)
   --help, -h                show help (default: false)

`)
}

func (s *serverProvider) stdLoggerPrint(format string, args ...any) {
	logger.Infof("[%s] %s", s.opt.AppName, fmt.Sprintf(format, args...))
}

func (s *serverProvider) stdErrLoggerPrint(format string, args ...any) {
	logger.Infof("[%s] %s", s.opt.AppName, fmt.Sprintf(format, args...))
}
