package router

import (
	"context"
	"mime/multipart"
	"net/http"
)

type Config struct {
	MaxMultipartMemory *int64
	UseRawPath         *bool
}

type HttpRouter interface {
	http.Handler
	IRouter

	SetConfig(conf Config)
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

	Next()
	Abort()

	Set(key string, value any)
	Get(key string) (any, bool)

	FullPath() string
	Request() *http.Request
	GetRawData() ([]byte, error)
	Param(key string) string
	Query(key string) string
	PostForm(key string) string
	PostFormArray(key string) []string
	FormFile(name string) (*multipart.FileHeader, error)
	Cookie(name string) (string, error)
	GetHeader(key string) string

	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite)
	Status(code int)
	Header(key, value string)
	Response(code int, contentType string, data []byte)
}

type Handler func(ctx Context)

type Middleware func(ctx Context)
