package rest

import (
	"net/http"

	"github.com/hyper-micro/hyper/rest/adapter"
	"github.com/hyper-micro/hyper/rest/adapter/gin"
	"github.com/hyper-micro/hyper/rest/router"
)

type Rest struct {
	adapter router.HttpRouter
}

func New(adapterIDs ...string) (*Rest, error) {
	var adapterID = gin.AdaptName
	if len(adapterIDs) > 0 {
		adapterID = adapterIDs[0]
	}
	routerAdapter, err := adapter.Get(adapterID)
	if err != nil {
		return nil, err
	}
	return &Rest{
		adapter: routerAdapter,
	}, nil
}

func (r *Rest) Group(path string, middlewares ...router.Middleware) router.IRouter {
	return r.adapter.Group(path, middlewares...)
}

func (r *Rest) Use(middlewares ...router.Middleware) router.IRoutes {
	return r.adapter.Use(middlewares...)
}

func (r *Rest) Get(path string, handler router.Handler) router.IRoutes {
	return r.adapter.Get(path, handler)
}

func (r *Rest) Post(path string, handler router.Handler) router.IRoutes {
	return r.adapter.Post(path, handler)
}

func (r *Rest) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.adapter.ServeHTTP(w, req)
}
