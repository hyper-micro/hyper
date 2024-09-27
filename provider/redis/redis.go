package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"runtime"

	"github.com/hyper-micro/hyper/config"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

type Provider interface {
	Into(instance ...string) *redis.Client
}

type redisProvider struct {
	clients map[string]*redis.Client
}

func NewProvider(conf config.Config) (Provider, func(), error) {
	var clients = make(map[string]*redis.Client)
	cfg := conf.GetStringMap("db.redis")
	for k, c := range cfg {
		m, ok := c.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("db.redis configuration format error")
		}

		var (
			host = cast.ToString(m["host"])
			port = cast.ToInt(m["port"])
		)

		cpuNum := runtime.NumCPU()
		if cpuNum < 1 {
			cpuNum = 1
		}
		poolSize := cpuNum * 10

		var onConnect = func(ctx context.Context, cn *redis.Conn) error {
			return nil
		}

		var tlsConfig = &tls.Config{
			InsecureSkipVerify: cast.ToBool(m["skipVerify"]),
		}
		if !cast.ToBool(m["tls"]) {
			tlsConfig = nil
		}

		rdb := redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d",
				host,
				port,
			),
			Password:              cast.ToString(m["password"]),
			DB:                    cast.ToInt(m["db"]),
			ReadTimeout:           cast.ToDuration(m["timeout"]),
			WriteTimeout:          cast.ToDuration(m["timeout"]),
			DialTimeout:           cast.ToDuration(m["timeout"]),
			TLSConfig:             tlsConfig,
			OnConnect:             onConnect,
			ConnMaxIdleTime:       cast.ToDuration(m["maxIdleTime"]),
			PoolSize:              poolSize,
			ContextTimeoutEnabled: true,
			ConnMaxLifetime:       cast.ToDuration(m["maxLifetime"]),
			MaxRetries:            cast.ToInt(m["maxRetries"]),
		})

		rdb.AddHook(Hook{host})

		clients[k] = rdb
	}

	provider := &redisProvider{clients}

	return provider, provider.cleanup, nil
}

func (p *redisProvider) Into(instance ...string) *redis.Client {
	key := "default"
	if len(instance) > 0 {
		key = instance[0]
	}
	client, ok := p.clients[key]
	if !ok {
		panic(fmt.Sprintf("redis instance '%s' not initialized", key))
	}
	return client
}

func (p *redisProvider) cleanup() {
	for _, client := range p.clients {
		_ = client.Close()
	}
}
