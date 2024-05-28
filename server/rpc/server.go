package rpc

import (
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Config struct {
	Addr           string
	WriteBufSize   int
	ReadBufSize    int
	MaxRecvMsgSize int
	MaxSendMsgSize int
	Reflection     bool
}

type Option struct {
	Config

	ServiceOpts []grpc.ServerOption
}

type Server struct {
	opt      Option
	srv      *grpc.Server
	handlers []HandlerFn
}

type HandlerFn func(srv *grpc.Server)

func New(opt Option) *Server {

	srvOpts := []grpc.ServerOption{
		grpc.WriteBufferSize(opt.WriteBufSize),
		grpc.ReadBufferSize(opt.ReadBufSize),
		grpc.MaxRecvMsgSize(opt.MaxRecvMsgSize),
		grpc.MaxSendMsgSize(opt.MaxSendMsgSize),
	}

	srvOpts = append(srvOpts, opt.ServiceOpts...)

	srv := &Server{
		opt: opt,
		srv: grpc.NewServer(srvOpts...),
	}

	return srv
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.opt.Addr)
	if err != nil {
		return err
	}

	for _, h := range s.handlers {
		h(s.srv)
	}

	return s.srv.Serve(l)
}

func (s *Server) Shutdown() error {
	s.srv.GracefulStop()

	return nil
}

func (s *Server) Handler(handler HandlerFn) {
	h := func(srv *grpc.Server) {
		if s.opt.Reflection {
			reflection.Register(srv)
		}

		handler(srv)
	}
	s.handlers = append(s.handlers, h)
}
