package web

import (
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hyper-micro/hyper/internal/json"
	"github.com/spf13/cast"
)

type Ctx interface {
	context.Context

	Request() *http.Request
	Abort()
	IsAbort() bool

	Set(key string, value any)
	Get(key string) (value any, exists bool)
	GetString(key string) (s string)
	GetBool(key string) (b bool)
	GetInt(key string) (i int)
	GetInt64(key string) (i64 int64)
	GetFloat64(key string) (f64 float64)

	Param(key string) string
	ParamBool(key string) bool
	ParamInt(key string) int
	ParamInt64(key string) int64
	ParamFloat64(key string) float64

	QueryArray(key string) (values []string, ok bool)
	Query(key string) (string, bool)
	QueryString(key string) string
	QueryBool(key string) bool
	QueryInt(key string) int
	QueryInt64(key string) int64
	QueryFloat64(key string) float64

	PostFormArray(key string) (values []string, ok bool)
	PostForm(key string) (string, bool)
	PostFormString(key string) string
	PostFormBool(key string) bool
	PostFormInt(key string) int
	PostFormInt64(key string) int64
	PostFormFloat64(key string) float64
	FormFileHeader(name string) (*multipart.FileHeader, error)

	JsonBinding(d any) error

	GetHeader(key string) string
	Header(key, value string)
	Cookie(name string) (string, error)
	SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite)

	Status(code int)
	ResponseWithStatus(code int, data []byte) error
	Response(data []byte) error
	Json(data any) error
	String(data string) error
}

type ctx struct {
	ctx        context.Context
	srv        *Server
	mu         *sync.RWMutex
	w          http.ResponseWriter
	r          *http.Request
	kv         map[string]any
	abort      bool
	queryCache url.Values
	formCache  url.Values
	status     bool
}

const requestCtxKey = "_hyper/contextKey"

var validate = validator.New(validator.WithRequiredStructEnabled())

func makeContext(srv *Server, w http.ResponseWriter, r *http.Request) *ctx {
	c := r.Context()
	cCtx, ok := c.Value(requestCtxKey).(*ctx)
	if !ok {
		cCtx = newContext(c, srv, w, r)
		*r = *r.WithContext(context.WithValue(c, requestCtxKey, cCtx))
	}
	return cCtx
}

func newContext(c context.Context, srv *Server, w http.ResponseWriter, r *http.Request) *ctx {
	return &ctx{
		ctx: c,
		srv: srv,
		mu:  new(sync.RWMutex),
		w:   w,
		r:   r,
		kv:  make(map[string]any),
	}
}

func (c *ctx) Request() *http.Request {
	return c.r
}

func (c *ctx) Abort() {
	c.abort = true
}

func (c *ctx) IsAbort() bool {
	return c.abort
}

/// kv

func (c *ctx) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.kv[key] = value
}

func (c *ctx) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.kv[key]
	return
}

func (c *ctx) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *ctx) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *ctx) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *ctx) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *ctx) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// Params

func (c *ctx) Param(key string) string {
	vars := mux.Vars(c.r)
	return vars[key]
}

func (c *ctx) ParamBool(key string) bool {
	return cast.ToBool(c.Param(key))
}

func (c *ctx) ParamInt(key string) int {
	return cast.ToInt(c.Param(key))
}

func (c *ctx) ParamInt64(key string) int64 {
	return cast.ToInt64(c.Param(key))
}

func (c *ctx) ParamFloat64(key string) float64 {
	return cast.ToFloat64(c.Param(key))
}

// Queries

func (c *ctx) initQueryCache() {
	if c.queryCache == nil {
		if c.r != nil {
			c.queryCache = c.r.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *ctx) QueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

func (c *ctx) Query(key string) (string, bool) {
	if values, ok := c.QueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *ctx) QueryString(key string) string {
	value, _ := c.Query(key)
	return value
}

func (c *ctx) QueryBool(key string) bool {
	value, _ := c.Query(key)
	return cast.ToBool(value)
}

func (c *ctx) QueryInt(key string) int {
	value, _ := c.Query(key)
	return cast.ToInt(value)
}

func (c *ctx) QueryInt64(key string) int64 {
	value, _ := c.Query(key)
	return cast.ToInt64(value)
}

func (c *ctx) QueryFloat64(key string) float64 {
	value, _ := c.Query(key)
	return cast.ToFloat64(value)
}

/// PostForm

func (c *ctx) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		_ = c.r.ParseMultipartForm(c.srv.MaxMultipartMemory)
		c.formCache = c.r.PostForm
	}
}

func (c *ctx) PostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

func (c *ctx) PostForm(key string) (string, bool) {
	if values, ok := c.PostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *ctx) PostFormString(key string) string {
	value, _ := c.PostForm(key)
	return value
}

func (c *ctx) PostFormBool(key string) bool {
	value, _ := c.PostForm(key)
	return cast.ToBool(value)
}

func (c *ctx) PostFormInt(key string) int {
	value, _ := c.PostForm(key)
	return cast.ToInt(value)
}

func (c *ctx) PostFormInt64(key string) int64 {
	value, _ := c.PostForm(key)
	return cast.ToInt64(value)
}

func (c *ctx) PostFormFloat64(key string) float64 {
	value, _ := c.PostForm(key)
	return cast.ToFloat64(value)
}

func (c *ctx) FormFileHeader(name string) (*multipart.FileHeader, error) {
	if c.r.MultipartForm == nil {
		if err := c.r.ParseMultipartForm(c.srv.MaxMultipartMemory); err != nil {
			return nil, err
		}
	}
	f, fh, err := c.r.FormFile(name)
	if err != nil {
		return nil, err
	}
	_ = f.Close()
	return fh, err
}

/// Header

func (c *ctx) GetHeader(key string) string {
	return c.r.Header.Get(key)
}

func (c *ctx) Header(key, value string) {
	if value == "" {
		c.w.Header().Del(key)
		return
	}
	c.w.Header().Set(key, value)
}

/// Binding struct

func (c *ctx) JsonBinding(d any) error {
	dc := json.NewDecoder(c.r.Body)
	if err := dc.Decode(d); err != nil {
		return err
	}
	return validate.Struct(d)
}

/// Cookie

func (c *ctx) Cookie(name string) (string, error) {
	cookie, err := c.r.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *ctx) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.w, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

/// Response

func (c *ctx) Status(code int) {
	if !c.status {
		c.w.WriteHeader(code)
	}
	c.status = true
}

func (c *ctx) ResponseWithStatus(code int, data []byte) error {
	c.Status(code)
	_, err := c.w.Write(data)
	return err
}

func (c *ctx) Response(data []byte) error {
	return c.ResponseWithStatus(http.StatusOK, data)
}

func (c *ctx) Json(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	return c.Response(b)
}

func (c *ctx) String(data string) error {
	c.Header("Content-Type", "text/plain; charset=utf-8")
	return c.Response([]byte(data))
}

/// context.Context

func (c *ctx) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *ctx) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *ctx) Err() error {
	return c.ctx.Err()
}

func (c *ctx) Value(key any) any {
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	return c.ctx.Value(key)
}
