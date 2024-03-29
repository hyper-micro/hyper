package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Handler func(ctx *Context)

type MiddlewareHandler func(ctx *Context, next func())

type router struct {
	r   *mux.Router
	srv *Server
}

func newRouter(srv *Server, r *mux.Router) *router {
	return &router{
		srv: srv,
		r:   r,
	}
}

func (r *router) wrapHandler(f Handler) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := makeContext(r.srv, w, req)
		if ctx.IsAbort() {
			return
		}
		f(ctx)
	}
}

func (r *router) wrapMiddleware(f MiddlewareHandler) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := makeContext(r.srv, w, req)
			if ctx.IsAbort() {
				return
			}
			f(ctx, func() {
				next.ServeHTTP(w, req)
			})
		})
	}
}

func (r *router) Get(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodGet)
}

func (r *router) Head(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodHead)
}

func (r *router) Post(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodPost)
}

func (r *router) Put(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodPut)
}

func (r *router) Patch(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodPatch)
}

func (r *router) Delete(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodDelete)
}

func (r *router) Connect(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodConnect)
}

func (r *router) Options(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodOptions)
}

func (r *router) Trace(path string, f Handler) {
	r.r.HandleFunc(path, r.wrapHandler(f)).Methods(http.MethodTrace)
}

func (r *router) Use(fs ...MiddlewareHandler) {
	var nfs []mux.MiddlewareFunc
	for _, f := range fs {
		nfs = append(nfs, r.wrapMiddleware(f))
	}
	r.r.Use(nfs...)
}

func (r *router) PathPrefix(prefix string) *router {
	nr := r.r.PathPrefix(prefix).Subrouter()
	return newRouter(r.srv, nr)
}

func (r *router) HostPrefix(host string) *router {
	nr := r.r.Host(host).Subrouter()
	return newRouter(r.srv, nr)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}
