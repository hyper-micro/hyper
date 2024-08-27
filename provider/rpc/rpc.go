package rpc

import (
	"math"

	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/server/rpc"
)

type Provider interface {
	Into() *rpc.Server
	Run() error
	Shutdown() error
	Addr() string
}

type rpcProvider struct {
	addr string
	opt  rpc.Option
	srv  *rpc.Server
	conf config.Config
}

type ServerOption func(option *rpc.Option)

func NewProvider(conf config.Config, serverOptions ...ServerOption) Provider {
	addr := conf.GetStringOrDefault("server.rpc.addr", "0.0.0.0:18110")
	opt := rpc.Option{
		Config: rpc.Config{
			Addr:           addr,
			WriteBufSize:   conf.GetIntOrDefault("server.rpc.writeBufSize", 32*1024),
			ReadBufSize:    conf.GetIntOrDefault("server.rpc.readBufSize", 32*1024),
			MaxRecvMsgSize: conf.GetIntOrDefault("server.rpc.maxRecvMsgSize", 1024*1024*4),
			MaxSendMsgSize: conf.GetIntOrDefault("server.rpc.maxSendMsgSize", math.MaxInt32),
			Reflection:     conf.GetBoolOrDefault("server.rpc.reflection", false),
		},
		ServiceOpts: nil,
	}

	for _, apply := range serverOptions {
		apply(&opt)
	}

	p := &rpcProvider{
		addr: addr,
		opt:  opt,
		srv:  rpc.New(opt),
		conf: conf,
	}

	return p
}

func (p *rpcProvider) Into() *rpc.Server {
	return p.srv
}

func (p *rpcProvider) Run() error {
	return p.srv.Run()
}

func (p *rpcProvider) Shutdown() error {
	return p.srv.Shutdown()
}

func (p *rpcProvider) Addr() string {
	return p.addr
}
