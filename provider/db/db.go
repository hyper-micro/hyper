package db

import (
	"fmt"
	"os"
	"slices"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hyper-micro/hyper/config"
	"github.com/spf13/cast"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

type Provider interface {
	Into(instance ...string) *xorm.Engine
}

type dbProvider struct {
	engines map[string]*xorm.Engine
}

func NewProvider(conf config.Config) (Provider, func(), error) {
	var engines = make(map[string]*xorm.Engine)

	cfg := conf.GetStringMap("db.db")
	for k, c := range cfg {
		m, ok := c.(map[string]interface{})
		if !ok {
			return nil, nil, fmt.Errorf("db.db configuration format error")
		}
		dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=%s",
			m["username"],
			m["password"],
			m["host"],
			m["port"],
			m["dbname"],
			m["charset"],
		)

		driver := cast.ToString(m["driver"])
		if !slices.Contains([]string{"mysql", "pg"}, driver) {
			return nil, nil, fmt.Errorf("db.db.%s.driver not supported, driver = {%s}", k, driver)
		}

		engine, err := xorm.NewEngine(cast.ToString(m["driver"]), dsn)
		if err != nil {
			return nil, nil, err
		}

		engine.SetMaxIdleConns(cast.ToInt(m["maxIdleConns"]))
		engine.SetMaxOpenConns(cast.ToInt(m["maxOpenConns"]))
		engine.SetConnMaxLifetime(cast.ToDuration(m["maxLifetime"]))

		engine.SetLogger(
			log.NewSimpleLogger3(os.Stdout, log.DEFAULT_LOG_PREFIX, log.DEFAULT_LOG_FLAG, log.LOG_WARNING),
		)

		engine.AddHook(Hook{
			host:     cast.ToString(m["host"]),
			port:     cast.ToInt(m["port"]),
			database: cast.ToString(m["dbname"]),
		})
		engines[k] = engine
	}

	provider := &dbProvider{engines}

	return provider, provider.cleanup, nil
}

func (p *dbProvider) Into(instance ...string) *xorm.Engine {
	key := "default"
	if len(instance) > 0 {
		key = instance[0]
	}
	engine, ok := p.engines[key]
	if !ok {
		panic(fmt.Sprintf("db instance '%s' not initialized", key))
	}
	return engine
}

func (p *dbProvider) cleanup() {
	for _, engine := range p.engines {
		_ = engine.Close()
	}
}
