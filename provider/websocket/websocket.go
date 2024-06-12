package websocket

import (
	"net/http"
	"time"

	"github.com/hyper-micro/hyper/config"
	"github.com/hyper-micro/hyper/provider/logger"
	"github.com/hyper-micro/hyper/server/websocket"
)

type Provider interface {
	Into() *websocket.Server
	Run() error
	Shutdown() error
	Addr() string
}

type websocketProvider struct {
	addr string
	srv  *websocket.Server
}

func NewProvider(conf config.Config, logger logger.Provider) Provider {
	addr := conf.GetStringOrDefault("server.websocket.addr", "0.0.0.0:18110")
	readTimeout := conf.GetDurationOrDefault("server.websocket.readTimeout", time.Second)
	readBuffer := conf.GetIntOrDefault("server.websocket.readBuffer", 32*1024)
	opt := websocket.Option{
		Config: websocket.Config{
			Addr:              addr,
			HandshakeTimeout:  conf.GetDurationOrDefault("server.websocket.handshakeTimeout", readTimeout),
			ReadBuffer:        readBuffer,
			WriteBuffer:       conf.GetIntOrDefault("server.websocket.writeBuffer", readBuffer),
			ReadTimeout:       readTimeout,
			ReadHeaderTimeout: conf.GetDurationOrDefault("server.websocket.readTimeout", readTimeout),
			WriteTimeout:      conf.GetDurationOrDefault("server.websocket.writeTimeout", readTimeout),
			ShutdownTimeout:   conf.GetDurationOrDefault("server.websocket.shutdownTimeout", 5*time.Second),
			CertFile:          conf.GetString("server.websocket.certFile"),
			KeyFile:           conf.GetString("server.websocket.keyFile"),
		},
		Logger: logger.Into(),
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		MessageType: websocket.BinaryMessage,
	}

	p := &websocketProvider{
		addr: addr,
		srv:  websocket.New(opt),
	}

	return p
}

func (p *websocketProvider) Into() *websocket.Server {
	return p.srv
}

func (p *websocketProvider) Run() error {
	return p.srv.Run()
}

func (p *websocketProvider) Shutdown() error {
	return p.srv.Shutdown()
}

func (p *websocketProvider) Addr() string {
	return p.addr
}
