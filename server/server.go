package server

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
)

type BehaveHandler func(Server) error

type Server interface {
	Name() string
	BeforeRunHandler() []BehaveHandler
	BeforeShutdownHandler() []BehaveHandler
	AfterStopHandler() []BehaveHandler
	Run() error
	Shutdown() error
}

type Option struct {
	ShutdownSignal  []os.Signal
	ShutdownHandler func(error)
}

type server struct {
	Option

	srvs     []Server
	shutting bool
	shutChan chan struct{}
	mu       *sync.Mutex
}

func New(opt Option, srvs ...Server) Server {
	return &server{
		Option:   opt,
		srvs:     srvs,
		mu:       new(sync.Mutex),
		shutChan: make(chan struct{}),
	}
}

func (s *server) Name() string {
	return "HyperRestServer"
}

func (s *server) Run() error {
	var (
		runErrChan = make(chan error, len(s.srvs))
	)

	for _, srv := range s.srvs {
		go func(srv Server) {
			for _, fn := range srv.BeforeRunHandler() {
				if err := fn(srv); err != nil {
					runErrChan <- fmt.Errorf("hyper.server: '%s' beforeRun err: %v", srv.Name(), err)
					return
				}
			}

			var errs []error
			if err := srv.Run(); err != nil {
				errs = append(
					errs,
					fmt.Errorf("hyper.server: '%s' run err: %v", srv.Name(), err),
				)
			}
			for _, fn := range srv.AfterStopHandler() {
				if err := fn(srv); err != nil {
					errs = append(
						errs,
						fmt.Errorf("hyper.server: '%s' afterStop err: %v", srv.Name(), err),
					)
				}
			}

			runErrChan <- wrapError(errs)
		}(srv)
	}

	if len(s.Option.ShutdownSignal) > 0 {
		go func() {
			shutdownSign := make(chan os.Signal, 1)
			signal.Notify(shutdownSign, s.Option.ShutdownSignal...)
			<-shutdownSign
			_ = s.shutdown()
		}()
	}

	var errs []error
	for range s.srvs {
		if err := <-runErrChan; err != nil {
			errs = append(errs, err)
			if shutdownErr := s.shutdown(); shutdownErr != nil {
				errs = append(errs, shutdownErr)
			}
		}
	}

	if shutdownErr := s.shutdown(); shutdownErr != nil {
		errs = append(errs, shutdownErr)
	}

	<-s.shutChan

	return wrapError(errs)
}

func (s *server) Shutdown() error {
	return s.shutdown()
}

func (s *server) shutdown() error {
	s.mu.Lock()
	shutting := s.shutting
	s.shutting = true
	s.mu.Unlock()
	if shutting {
		return nil
	}

	var errs []error
	for _, srv := range s.srvs {
		for _, fn := range srv.BeforeShutdownHandler() {
			if err := fn(srv); err != nil {
				errs = append(
					errs,
					fmt.Errorf("hyper.server: '%s' beforeShutdown err: %v", srv.Name(), err),
				)
			}
		}

		if err := srv.Shutdown(); err != nil {
			errs = append(
				errs,
				fmt.Errorf("hyper.server: '%s' shutdown err: %v", srv.Name(), err),
			)
		}
	}

	wrapErr := wrapError(errs)
	if s.Option.ShutdownHandler != nil {
		s.Option.ShutdownHandler(wrapErr)
	}

	close(s.shutChan)

	return wrapErr
}

func (s *server) BeforeRunHandler() []BehaveHandler      { panic("should not be reached") }
func (s *server) BeforeShutdownHandler() []BehaveHandler { panic("should not be reached") }
func (s *server) AfterStopHandler() []BehaveHandler      { panic("should not be reached") }

func wrapError(errs []error) error {
	var wrapErr error
	for _, err := range errs {
		if wrapErr == nil {
			wrapErr = err
		} else {
			wrapErr = fmt.Errorf("%v: %w", err, wrapErr)
		}
	}
	return wrapErr
}
