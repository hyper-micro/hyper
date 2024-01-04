package rpc

import (
	"net"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/engine"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Name     string
	Addr     string
	CertFile string
	KeyFile  string
}

type Option struct {
	Config

	BeforeRunHandler, BeforeShutdownHandler, AfterStopHandler engine.BehaveHandler
	ServerOption                                              []grpc.ServerOption
	Handler                                                   func(srv *grpc.Server)
	Reflection                                                bool
	Prometheus                                                bool
}

type Server struct {
	Option

	srv *grpc.Server
}

func NewServer(opt Option) engine.Server {
	if opt.Prometheus {
		opt.ServerOption = append(
			opt.ServerOption,
			grpc.ChainUnaryInterceptor(
				grpc_prometheus.UnaryServerInterceptor,
			),
			grpc.StreamInterceptor(
				grpc_prometheus.StreamServerInterceptor,
			),
		)
	}

	return &Server{
		Option: opt,
	}
}

func NewDefaultServer(handler func(srv *grpc.Server), behaveHandler ...engine.BehaveHandlerSet) engine.Server {
	opt := Option{
		Config: Config{
			Name:     config.Default().GetStringOrDefault("server.rpc.name", "RpcServer"),
			Addr:     config.Default().GetStringOrDefault("server.rpc.addr", ":8081"),
			CertFile: config.Default().GetString("server.rpc.certFile"),
			KeyFile:  config.Default().GetString("server.rpc.keyFile"),
		},
		Handler: handler,
	}
	if len(behaveHandler) > 0 {
		if behaveHandler[0].BeforeHandler != nil {
			opt.BeforeRunHandler = behaveHandler[0].BeforeHandler
		}
		if behaveHandler[0].BeforeShutdownHandler != nil {
			opt.BeforeShutdownHandler = behaveHandler[0].BeforeShutdownHandler
		}
		if behaveHandler[0].AfterHandler != nil {
			opt.AfterStopHandler = behaveHandler[0].AfterHandler
		}
	}
	rpcSrv := NewServer(opt)

	return rpcSrv
}

func (s *Server) Name() string {
	return "GRPCServer"
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
	lis, err := net.Listen("tcp", s.Option.Config.Addr)
	if err != nil {
		return err
	}

	if s.Option.Config.CertFile != "" && s.Option.Config.KeyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(s.Option.Config.CertFile, s.Option.Config.KeyFile)
		if err != nil {
			return err
		}
		s.Option.ServerOption = append(s.Option.ServerOption, grpc.Creds(creds))
	}

	s.srv = grpc.NewServer(s.Option.ServerOption...)

	if s.Option.Handler != nil {
		s.Option.Handler(s.srv)
	}
	if s.Option.Reflection {
		reflection.Register(s.srv)
	}
	if s.Option.Prometheus {
		grpc_prometheus.EnableHandlingTimeHistogram()
		grpc_prometheus.Register(s.srv)
	}

	return s.srv.Serve(lis)
}

func (s *Server) Shutdown() error {
	if s.srv != nil {
		s.srv.GracefulStop()
	}
	return nil
}
