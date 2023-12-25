package router

import (
	"context"
	"net/http"
)

type HttpRouter interface {
	http.Handler
	IRouter
}

type IRouter interface {
	IRoutes
	Group(string, ...Middleware) IRouter
}

type IRoutes interface {
	Use(...Middleware) IRoutes
	Get(string, Handler) IRoutes
	Post(string, Handler) IRoutes
}

type Context interface {
	context.Context
	GetRawData() ([]byte, error)
	Next()
	Abort()
	FullPath() string
	Set(key string, value any)
	Get(key string) (any, bool)
	Param(key string) string
	Query(key string) string
}

type Handler func(ctx Context)

type Middleware func(ctx Context)
