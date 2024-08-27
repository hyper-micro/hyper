package websocket

import (
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hyper-micro/hyper/logger"
	"github.com/hyper-micro/hyper/toolkit/slice"
)

type Config struct {
	Addr                           string
	HandshakeTimeout               time.Duration
	ReadBuffer                     int
	WriteBuffer                    int
	ReadTimeout, ReadHeaderTimeout time.Duration
	WriteTimeout                   time.Duration
	ShutdownTimeout                time.Duration
	CertFile, KeyFile              string
}

type Handler interface {
	OnConnection(rwc net.Conn)
}

type Option struct {
	Config

	Logger         logger.Logger
	MaxHeaderBytes int
	ConnState      func(net.Conn, http.ConnState)
	TLSConfig      *tls.Config
	TLSNextProto   map[string]func(*http.Server, *tls.Conn, http.Handler)
	ErrorLog       *log.Logger
	BaseContext    func(net.Listener) context.Context
	ConnContext    func(ctx context.Context, c net.Conn) context.Context
	CheckOrigin    func(r *http.Request) bool
	MessageType    int
}

type Server struct {
	Option

	up      websocket.Upgrader
	srv     *http.Server
	handler Handler
}

func New(opt Option) *Server {
	ws := &Server{
		Option: opt,
		up: websocket.Upgrader{
			ReadBufferSize:   opt.ReadBuffer,
			WriteBufferSize:  opt.WriteBuffer,
			HandshakeTimeout: opt.HandshakeTimeout,
			CheckOrigin:      opt.CheckOrigin,
		},
	}

	ws.srv = &http.Server{
		Addr:              opt.Addr,
		Handler:           ws,
		TLSConfig:         opt.TLSConfig,
		ReadTimeout:       opt.ReadTimeout,
		ReadHeaderTimeout: opt.ReadHeaderTimeout,
		WriteTimeout:      opt.WriteTimeout,
		MaxHeaderBytes:    opt.MaxHeaderBytes,
		TLSNextProto:      opt.TLSNextProto,
		ConnState:         opt.ConnState,
		ErrorLog:          opt.ErrorLog,
		BaseContext:       opt.BaseContext,
		ConnContext:       opt.ConnContext,
	}
	return ws
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

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !slice.ContainsString("mqtt", websocket.Subprotocols(r)) {
		s.Option.Logger.Errorf("websocket: Client does not support mqtt sub Protocol, remoteAddr: %s", r.RemoteAddr)
		return
	}

	var header = make(http.Header)
	header.Add("Sec-Websocket-Protocol", "mqtt")
	wsConn, err := s.up.Upgrade(w, r, header)
	if err != nil {
		s.Option.Logger.Errorf("websocket: Connect err: %s, remoteAddr: %s", err.Error(), r.RemoteAddr)
		return
	}

	rwc := newConn(s, wsConn)

	s.Option.Logger.Infof("websocket: Connection established, localAddr: %s, remoteAddr: %s", rwc.LocalAddr().String(), rwc.RemoteAddr().String())

	if s.handler != nil {
		s.handler.OnConnection(rwc)
	}
}

func (s *Server) Handler(handler Handler) {
	s.handler = handler
}
