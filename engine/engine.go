package engine

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/hyper-micro/hyper/errors"
	"github.com/hyper-micro/hyper/logger"
)

type BehaveHandler func(engine *Engine, srv Server) error

type BehaveHandlerSet struct {
	BeforeHandler, AfterHandler, BeforeShutdownHandler BehaveHandler
}

type Server interface {
	Name() string
	BeforeRunHandler() BehaveHandler
	BeforeShutdownHandler() BehaveHandler
	AfterStopHandler() BehaveHandler
	Run() error
	Shutdown() error
}

type Option struct {
	Name            string
	ShutdownSignal  []os.Signal
	ShutdownHandler func(error)
	Console         bool
	BeforeHandler   BehaveHandler
	AfterHandler    BehaveHandler
}

type Engine struct {
	Option

	running  bool
	shutting bool
	servers  []Server
	shutChan chan struct{}
	mu       *sync.Mutex
}

func New(opt Option, servers ...Server) *Engine {
	return &Engine{
		Option:   opt,
		servers:  servers,
		mu:       new(sync.Mutex),
		shutChan: make(chan struct{}),
	}
}

func (s *Engine) Name() string {
	return s.Option.Name
}

func (s *Engine) Run() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}
	s.mu.Unlock()
	s.running = true

	s.Print("[%s] server boot, pid is %d", s.Name(), os.Getpid())

	defer func() {
		s.Print("[%s] server stop, Bye!", s.Name())
	}()

	if s.Option.BeforeHandler != nil {
		if err := s.Option.BeforeHandler(s, s); err != nil {
			return fmt.Errorf("engine: beforeHandler err: %v", err)
		}
	}

	var (
		runErrChan = make(chan error, len(s.servers))
	)

	for _, srv := range s.servers {
		go func(srv Server) {
			if fn := srv.BeforeRunHandler(); fn != nil {
				if err := fn(s, srv); err != nil {
					runErrChan <- fmt.Errorf("hyper.server: '%s' beforeRun err: %v", srv.Name(), err)
					return
				}
			}

			var wrapErr error

			s.Print("[%s] server start", srv.Name())

			if err := srv.Run(); err != nil {
				wrapErr = errors.Wrap(wrapErr, fmt.Errorf("hyper.server: '%s' run err: %v", srv.Name(), err))
			}

			if fn := srv.AfterStopHandler(); fn != nil {
				if err := fn(s, srv); err != nil {
					wrapErr = errors.Wrap(wrapErr, fmt.Errorf("hyper.server: '%s' afterStop err: %v", srv.Name(), err))
				}
			}

			runErrChan <- wrapErr
		}(srv)
	}

	if len(s.Option.ShutdownSignal) > 0 {
		go func() {
			shutdownSign := make(chan os.Signal, 1)
			signal.Notify(shutdownSign, s.Option.ShutdownSignal...)

			s.Print("[%s] signal.Notify: %v", s.Name(), s.Option.ShutdownSignal)

			recSign := <-shutdownSign

			s.Print("[%s] receive signal: %v", s.Name(), recSign)

			_ = s.shutdown()
		}()
	}

	var wrapErr error
	for range s.servers {
		if err := <-runErrChan; err != nil {
			wrapErr = errors.Wrap(wrapErr, err)
			if shutdownErr := s.shutdown(); shutdownErr != nil {
				wrapErr = errors.Wrap(wrapErr, err)
			}
		}
	}

	if shutdownErr := s.shutdown(); shutdownErr != nil {
		wrapErr = errors.Wrap(wrapErr, shutdownErr)
	}

	<-s.shutChan

	if s.Option.AfterHandler != nil {
		if err := s.Option.AfterHandler(s, s); err != nil {
			wrapErr = errors.Wrap(wrapErr, fmt.Errorf("engine: afterHandler err: %v", err))
		}
	}

	return wrapErr
}

func (s *Engine) Shutdown() error {
	return s.shutdown()
}

func (s *Engine) shutdown() error {
	s.mu.Lock()
	if s.shutting {
		s.mu.Unlock()
		return nil
	}
	s.mu.Unlock()
	s.shutting = true

	var wrapErr error
	for _, srv := range s.servers {
		if fn := srv.BeforeShutdownHandler(); fn != nil {
			s.Print("[%s] run beforeShutdownHandler: '%#v'", srv.Name(), fn)

			if err := fn(s, srv); err != nil {
				wrapErr = errors.Wrap(wrapErr, fmt.Errorf("hyper.server: '%s' beforeShutdown err: %v", srv.Name(), err))
			}
		}

		s.Print("[%s] shutting down", srv.Name())

		if err := srv.Shutdown(); err != nil {
			wrapErr = errors.Wrap(wrapErr, fmt.Errorf("hyper.server: '%s' shutdown err: %v", srv.Name(), err))
		}
	}

	if s.Option.ShutdownHandler != nil {
		s.Option.ShutdownHandler(wrapErr)
	}

	close(s.shutChan)

	return wrapErr
}

func (s *Engine) Print(format string, args ...interface{}) {
	if s.Option.Console {
		logger.Infof(format, args...)
	}
}

func (s *Engine) BeforeRunHandler() BehaveHandler      { panic("unreachable") }
func (s *Engine) BeforeShutdownHandler() BehaveHandler { panic("unreachable") }
func (s *Engine) AfterStopHandler() BehaveHandler      { panic("unreachable") }
