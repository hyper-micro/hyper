package web

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Config struct {
	Addr                           string
	ReadTimeout, ReadHeaderTimeout time.Duration
	WriteTimeout, IdleTimeout      time.Duration
	ShutdownTimeout                time.Duration
	CertFile, KeyFile              string
}

type Option struct {
	Config

	MaxHeaderBytes     int
	MaxMultipartMemory int64
	ConnState          func(net.Conn, http.ConnState)
	TLSConfig          *tls.Config
	TLSNextProto       map[string]func(*http.Server, *tls.Conn, http.Handler)
	ErrorLog           *log.Logger
	BaseContext        func(net.Listener) context.Context
	ConnContext        func(ctx context.Context, c net.Conn) context.Context
}

type Server struct {
	Option

	*router
	srv      *http.Server
	handlers []HandlerFunc
}

func New(opt Option) *Server {
	srv := &Server{Option: opt}
	rr := mux.NewRouter()
	srv.router = newRouter(srv, rr)
	srv.srv = &http.Server{
		Addr:              opt.Addr,
		Handler:           srv.router,
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
	}

	return srv
}

func (s *Server) Run() error {
	if len(s.handlers) > 0 {
		h := s.srv.Handler
		for i := len(s.handlers) - 1; i >= 0; i-- {
			hf := s.handlers[i]
			h = &handler{
				h: h,
				f: hf,
			}
		}
		s.srv.Handler = h
	}

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

type HandlerFunc func(http.Handler, http.ResponseWriter, *http.Request)

func (s *Server) Handler(hf HandlerFunc) {
	s.handlers = append(s.handlers, hf)
}

type handler struct {
	h http.Handler
	f func(http.Handler, http.ResponseWriter, *http.Request)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.f(h.h, w, req)
}
