package adapter

import (
	"fmt"
	"sync"

	"github.com/hyper-micro/hyper/rest/router"
)

var (
	adapters = new(sync.Map)
)

func Register(id string, router router.HttpRouter) {
	adapters.Store(id, router)
}

func Get(id string) (router.HttpRouter, error) {
	val, ok := adapters.Load(id)
	if !ok {
		return nil, fmt.Errorf("rest: router adapter not found, adapterID: %v", id)
	}

	return val.(router.HttpRouter), nil
}
