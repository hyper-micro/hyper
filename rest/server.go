package rest

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/hyper-micro/hyper/engine"
)

type Config struct {
	Name                           string
	Addr                           string
	ReadTimeout, ReadHeaderTimeout time.Duration
	WriteTimeout, IdleTimeout      time.Duration
	ShutdownTimeout                time.Duration
	CertFile, KeyFile              string
}

type Option struct {
	Config

	MaxHeaderBytes                                            int
	MaxMultipartMemory                                        int64
	BeforeRunHandler, BeforeShutdownHandler, AfterStopHandler engine.BehaveHandler
	ConnState                                                 func(net.Conn, http.ConnState)
	TLSConfig                                                 *tls.Config
	TLSNextProto                                              map[string]func(*http.Server, *tls.Conn, http.Handler)
	ErrorLog                                                  *log.Logger
	BaseContext                                               func(net.Listener) context.Context
	ConnContext                                               func(ctx context.Context, c net.Conn) context.Context
}

type Server struct {
	Option

	srv *http.Server
}

func NewServer(opt Option, handler http.Handler) engine.Server {
	return &Server{
		Option: opt,
		srv: &http.Server{
			Addr:              opt.Addr,
			Handler:           handler,
			TLSConfig:         opt.TLSConfig,
			ReadTimeout:       opt.ReadTimeout,
			ReadHeaderTimeout: opt.ReadHeaderTimeout,
			WriteTimeout:      opt.WriteTimeout,
			IdleTimeout:       opt.IdleTimeout,
			MaxHeaderBytes:    opt.MaxHeaderBytes,
			TLSNextProto:      opt.TLSNextProto,
			ConnState:         opt.ConnState,
			ErrorLog:          opt.ErrorLog,
			BaseContext:       opt.BaseContext,
			ConnContext:       opt.ConnContext,
		},
	}
}

func (s *Server) Name() string {
	return s.Option.Name
}

func (s *Server) BeforeRunHandler() engine.BehaveHandler {
	return s.Option.BeforeRunHandler
}

func (s *Server) BeforeShutdownHandler() engine.BehaveHandler {
	return s.Option.BeforeShutdownHandler
}

func (s *Server) AfterStopHandler() engine.BehaveHandler {
	return s.Option.AfterStopHandler
}

func (s *Server) Run() error {
	var (
		srvErr error
	)
	for {
		if s.Option.CertFile != "" && s.Option.KeyFile != "" {
			srvErr = s.srv.ListenAndServeTLS(s.Option.CertFile, s.Option.KeyFile)
			break
		}
		srvErr = s.srv.ListenAndServe()
		break
	}
	if srvErr == http.ErrServerClosed {
		srvErr = nil
	}
	return srvErr
}

func (s *Server) Shutdown() error {
	shutdownTimeout := s.Option.ShutdownTimeout
	if shutdownTimeout == 0 {
		shutdownTimeout = 5 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return s.srv.Shutdown(ctx)
}
