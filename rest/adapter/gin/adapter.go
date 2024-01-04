package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyper-micro/hyper/rest/adapter"
	"github.com/hyper-micro/hyper/rest/router"
)

const AdaptName = "gin"

func init() {
	gin.SetMode(gin.ReleaseMode)
	adapter.Register(AdaptName, NewGinAdapter(gin.New()))
}

type Adapter struct {
	ginEngine *gin.Engine
	ginRouter gin.IRouter
	ginRoutes gin.IRoutes
}

func NewGinAdapter(engine *gin.Engine) router.HttpRouter {
	engine.RedirectTrailingSlash = false

	return &Adapter{
		ginEngine: engine,
		ginRouter: engine,
		ginRoutes: engine,
	}
}

func (g *Adapter) SetConfig(conf router.Config) {
	if conf.MaxMultipartMemory != nil {
		g.ginEngine.MaxMultipartMemory = *conf.MaxMultipartMemory
	}
	if conf.UseRawPath != nil {
		g.ginEngine.UseRawPath = *conf.UseRawPath
	}
}

func (g *Adapter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	g.ginEngine.ServeHTTP(w, req)
}

func (g *Adapter) HTTP404Handler() {
}

func (g *Adapter) Group(path string, middlewares ...router.Middleware) router.IRouter {
	var handlers []gin.HandlerFunc
	for _, h := range middlewares {
		handlers = append(handlers, func(ctx *gin.Context) {
			h(newGinAdapterContext(ctx))
		})
	}
	group := g.ginRouter.Group(path, handlers...)
	return &Adapter{
		ginRouter: group,
		ginRoutes: group,
	}
}

func (g *Adapter) Use(middlewares ...router.Middleware) router.IRoutes {
	var handlers []gin.HandlerFunc
	for _, h := range middlewares {
		handlers = append(handlers, func(ctx *gin.Context) {
			h(newGinAdapterContext(ctx))
		})
	}
	group := g.ginRouter.Use(handlers...)
	return &Adapter{
		ginRoutes: group,
	}
}

func (g *Adapter) Post(path string, handler router.Handler) router.IRoutes {
	r := g.ginRoutes.POST(path, func(ctx *gin.Context) {
		handler(newGinAdapterContext(ctx))
	})
	return &Adapter{
		ginRoutes: r,
	}
}

func (g *Adapter) Get(path string, handler router.Handler) router.IRoutes {
	r := g.ginRoutes.GET(path, func(ctx *gin.Context) {
		handler(newGinAdapterContext(ctx))
	})
	return &Adapter{
		ginRoutes: r,
	}
}
