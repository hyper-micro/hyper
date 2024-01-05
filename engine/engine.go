package engine

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/hyper-micro/hyper/errors"
	"github.com/hyper-micro/hyper/log"
)

type BehaveHandler func(Server) error

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

type engine struct {
	Option

	servers  []Server
	shutting bool
	shutChan chan struct{}
	mu       *sync.Mutex
}

func New(opt Option, servers ...Server) Server {
	return &engine{
		Option:   opt,
		servers:  servers,
		mu:       new(sync.Mutex),
		shutChan: make(chan struct{}),
	}
}

func (s *engine) Name() string {
	return s.Option.Name
}

func (s *engine) Run() error {

	s.Print("[%s] server boot, pid is %d", s.Name(), os.Getpid())

	defer func() {
		s.Print("[%s] server stop, Bye!", s.Name())
	}()

	if s.Option.BeforeHandler != nil {
		if err := s.Option.BeforeHandler(s); err != nil {
			return fmt.Errorf("engine: beforeHandler err: %v", err)
		}
	}

	var (
		runErrChan = make(chan error, len(s.servers))
	)

	for _, srv := range s.servers {
		go func(srv Server) {
			if fn := srv.BeforeRunHandler(); fn != nil {
				s.Print("[%s] run beforeRunHandler: '%#v'", srv.Name(), fn)

				if err := fn(srv); err != nil {
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
				s.Print("[%s] run afterStopHandler: '%#v'", srv.Name(), fn)

				if err := fn(srv); err != nil {
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
		if err := s.Option.AfterHandler(s); err != nil {
			wrapErr = errors.Wrap(wrapErr, fmt.Errorf("engine: afterHandler err: %v", err))
		}
	}

	return wrapErr
}

func (s *engine) Shutdown() error {
	return s.shutdown()
}

func (s *engine) shutdown() error {
	s.mu.Lock()
	shutting := s.shutting
	s.shutting = true
	s.mu.Unlock()
	if shutting {
		return nil
	}

	var wrapErr error
	for _, srv := range s.servers {
		if fn := srv.BeforeShutdownHandler(); fn != nil {
			s.Print("[%s] run beforeShutdownHandler: '%#v'", srv.Name(), fn)

			if err := fn(srv); err != nil {
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

func (s *engine) Print(format string, args ...interface{}) {
	if s.Option.Console {
		log.Infof(format, args...)
	}
}

func (s *engine) BeforeRunHandler() BehaveHandler      { panic("unreachable") }
func (s *engine) BeforeShutdownHandler() BehaveHandler { panic("unreachable") }
func (s *engine) AfterStopHandler() BehaveHandler      { panic("unreachable") }
