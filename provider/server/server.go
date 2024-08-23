package server

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/errors"
	"github.com/hyper-micro/hyper/logger"
)

type Provider interface {
	RegInit(fs ...RegInitHandler)
	RegServes(fs ...RegServeHandler) error
	Run() error
}

type App interface {
	Run() error
	Shutdown() error
	Addr() string
	Name() string
}

type CleanUpHandler func()

type RegServeHandler func(config.Config) (App, CleanUpHandler, error)

type RegInitHandler func(config.Config) error

type serverProvider struct {
	opt             Option
	apps            []App
	inits           []RegInitHandler
	cleanUps        []func()
	flagSet         *flag.FlagSet
	configFileFlag  string
	showHelpFlag    bool
	showVersionFlag bool
	inShutdown      bool
	conf            config.Config
}

type Option struct {
	AppName               string
	AppDesc               string
	Version               string
	BuildCommit           string
	BuildDate             string
	ShutdownSigs          []os.Signal
	ShutdownDelayDuration time.Duration
	ConfigPathType        config.PathType
	ConfigIgnoreFileName  bool
	ConfigDefault         string
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

func (s *serverProvider) RegServes(fs ...RegServeHandler) error {
	for _, f := range fs {
		app, cleanUp, err := f(s.conf)
		if err != nil {
			return err
		}
		if app != nil {
			s.apps = append(s.apps, app)
		}
		if cleanUp != nil {
			s.cleanUps = append(s.cleanUps, cleanUp)
		}
	}

	return nil
}

func (s *serverProvider) RegInit(fs ...RegInitHandler) {
	s.inits = append(s.inits, fs...)
}

func (s *serverProvider) Run() error {
	if len(s.opt.ShutdownSigs) > 0 {
		shutdownSignChan := make(chan os.Signal, 1)
		signal.Notify(shutdownSignChan, s.opt.ShutdownSigs...)
		go func() {
			recSign := <-shutdownSignChan
			s.stdLoggerPrint("Receive signal: %v", recSign)

			if s.opt.ShutdownDelayDuration > 0 {
				s.stdLoggerPrint("ShutdownDelayDuration %s", s.opt.ShutdownDelayDuration.String())
				time.Sleep(s.opt.ShutdownDelayDuration)
			}

			s.shutdown()
		}()
	}

	defer func() {
		s.stdLoggerPrint("Server stopped, Bye!")
	}()

	s.stdLoggerPrint("Load config file: %v", s.conf.FileNames())
	s.stdLoggerPrint("Version: %s, Commit: %s, buildDate: %s", s.opt.Version, s.opt.BuildCommit, s.opt.BuildDate)
	s.stdLoggerPrint("Pid: %v", os.Getpid())
	s.stdLoggerPrint("Signal.Notify: %v", s.opt.ShutdownSigs)

	for _, init := range s.inits {
		if err := init(s.conf); err != nil {
			return err
		}
	}

	var (
		appErr error
		wg     sync.WaitGroup
	)
	for _, app := range s.apps {
		wg.Add(1)
		go func(app App) {
			defer wg.Done()

			s.stdLoggerPrint("%s listen: %s", app.Name(), app.Addr())
			if err := app.Run(); err != nil {
				s.stdErrLoggerPrint("%s run error: %v", app.Name(), err)
				appErr = errors.Wrap(appErr, err)
			}

			s.shutdown()
		}(app)
	}

	wg.Wait()

	return appErr
}

func (s *serverProvider) shutdown() {
	if s.inShutdown {
		return
	}

	s.inShutdown = true

	for _, app := range s.apps {
		s.stdLoggerPrint("%s shutting down", app.Name())
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
	s.flagSet.StringVar(&s.configFileFlag, "c", s.opt.ConfigDefault, "set configure file path")
	s.flagSet.StringVar(&s.configFileFlag, "config", s.opt.ConfigDefault, "set configure file path")
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

	conf, err := config.New(s.opt.ConfigPathType, s.opt.ConfigIgnoreFileName, s.configFileFlag)
	if err != nil {
		return err
	}
	s.conf = conf

	return nil
}

func (s *serverProvider) flagUsage() {
	fmt.Printf(`
USAGE:
   %s [options]

%s

OPTIONS:
   --config value, -c value  set configure file path (default: "%s")
   --version, -v             show version (default: false)
   --help, -h                show help (default: false)

`, s.opt.AppName, s.opt.AppDesc, s.opt.ConfigDefault)
}

func (s *serverProvider) stdLoggerPrint(format string, args ...any) {
	logger.Infof("[%s] %s", s.opt.AppName, fmt.Sprintf(format, args...))
}

func (s *serverProvider) stdErrLoggerPrint(format string, args ...any) {
	logger.Errorf("[%s] %s", s.opt.AppName, fmt.Sprintf(format, args...))
}
