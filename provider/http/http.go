package http

import (
	"time"

	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/server/web"
)

type Provider interface {
	Into() *web.Server
	Run() error
	Shutdown() error
	Addr() string
}

type httpProvider struct {
	addr string
	opt  web.Option
	srv  *web.Server
	conf config.Config
}

func NewProvider(conf config.Config) Provider {
	addr := conf.GetStringOrDefault("server.http.addr", ":8080")
	timeout := conf.GetDurationOrDefault("server.http.timeout", 30*time.Second)
	opt := web.Option{
		Config: web.Config{
			Addr:            addr,
			ReadTimeout:     conf.GetDurationOrDefault("server.http.readTimeout", timeout),
			WriteTimeout:    conf.GetDurationOrDefault("server.http.writeTimeout", timeout),
			ShutdownTimeout: conf.GetDurationOrDefault("server.http.shutdown", 5*time.Second),
			CertFile:        conf.GetString("server.http.certFile"),
			KeyFile:         conf.GetString("server.http.keyFile"),
		},
	}

	p := &httpProvider{
		addr: addr,
		opt:  opt,
		srv:  web.New(opt),
		conf: conf,
	}

	return p
}

func (p *httpProvider) Into() *web.Server {
	return p.srv
}

func (p *httpProvider) Run() error {
	return p.srv.Run()
}

func (p *httpProvider) Shutdown() error {
	return p.srv.Shutdown()
}

func (p *httpProvider) Addr() string {
	return p.addr
}
