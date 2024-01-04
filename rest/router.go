package rest

import (
	"net/http"

	"github.com/hyper-micro/hyper/rest/adapter"
	"github.com/hyper-micro/hyper/rest/adapter/gin"
	"github.com/hyper-micro/hyper/rest/router"
)

type Router struct {
	adapter router.HttpRouter
}

func NewRouter(conf router.Config, adapterIDs ...string) (*Router, error) {
	var adapterID = gin.AdaptName
	if len(adapterIDs) > 0 {
		adapterID = adapterIDs[0]
	}
	routerAdapter, err := adapter.Get(adapterID)
	if err != nil {
		return nil, err
	}

	routerAdapter.SetConfig(conf)

	return &Router{
		adapter: routerAdapter,
	}, nil
}

func (r *Router) Group(path string, middlewares ...router.Middleware) router.IRouter {
	return r.adapter.Group(path, middlewares...)
}

func (r *Router) Use(middlewares ...router.Middleware) router.IRoutes {
	return r.adapter.Use(middlewares...)
}

func (r *Router) Get(path string, handler router.Handler) router.IRoutes {
	return r.adapter.Get(path, handler)
}

func (r *Router) Post(path string, handler router.Handler) router.IRoutes {
	return r.adapter.Post(path, handler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.adapter.ServeHTTP(w, req)
}
