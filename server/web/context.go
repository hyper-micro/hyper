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

type Context struct {
	ctx        context.Context
	srv        *Server
	mu         *sync.RWMutex
	w          http.ResponseWriter
	r          *http.Request
	kv         map[string]any
	abort      bool
	queryCache url.Values
	formCache  url.Values
}

const requestCtxKey = "_hyper/contextKey"

var validate = validator.New(validator.WithRequiredStructEnabled())

func makeContext(srv *Server, w http.ResponseWriter, r *http.Request) *Context {
	ctx := r.Context()
	cCtx, ok := ctx.Value(requestCtxKey).(*Context)
	if !ok {
		cCtx = newContext(ctx, srv, w, r)
		r.WithContext(context.WithValue(ctx, requestCtxKey, cCtx))
	}
	return cCtx
}

func newContext(ctx context.Context, srv *Server, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ctx: ctx,
		srv: srv,
		mu:  new(sync.RWMutex),
		w:   w,
		r:   r,
		kv:  make(map[string]any),
	}
}

func (c *Context) Request() *http.Request {
	return c.r
}

func (c *Context) Abort() {
	c.abort = true
}

func (c *Context) IsAbort() bool {
	return c.abort
}

/// kv

func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.kv[key] = value
}

func (c *Context) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, exists = c.kv[key]
	return
}

func (c *Context) GetString(key string) (s string) {
	if val, ok := c.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (c *Context) GetBool(key string) (b bool) {
	if val, ok := c.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (c *Context) GetInt(key string) (i int) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

func (c *Context) GetInt64(key string) (i64 int64) {
	if val, ok := c.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}

func (c *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := c.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}

// Params

func (c *Context) Param(key string) string {
	vars := mux.Vars(c.r)
	return vars[key]
}

func (c *Context) ParamBool(key string) bool {
	return cast.ToBool(c.Param(key))
}

func (c *Context) ParamInt(key string) int {
	return cast.ToInt(c.Param(key))
}

func (c *Context) ParamInt64(key string) int64 {
	return cast.ToInt64(c.Param(key))
}

func (c *Context) ParamFloat64(key string) float64 {
	return cast.ToFloat64(c.Param(key))
}

// Queries

func (c *Context) initQueryCache() {
	if c.queryCache == nil {
		if c.r != nil {
			c.queryCache = c.r.URL.Query()
		} else {
			c.queryCache = url.Values{}
		}
	}
}

func (c *Context) QueryArray(key string) (values []string, ok bool) {
	c.initQueryCache()
	values, ok = c.queryCache[key]
	return
}

func (c *Context) Query(key string) (string, bool) {
	if values, ok := c.QueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) QueryString(key string) string {
	value, _ := c.Query(key)
	return value
}

func (c *Context) QueryBool(key string) bool {
	value, _ := c.Query(key)
	return cast.ToBool(value)
}

func (c *Context) QueryInt(key string) int {
	value, _ := c.Query(key)
	return cast.ToInt(value)
}

func (c *Context) QueryInt64(key string) int64 {
	value, _ := c.Query(key)
	return cast.ToInt64(value)
}

func (c *Context) QueryFloat64(key string) float64 {
	value, _ := c.Query(key)
	return cast.ToFloat64(value)
}

/// PostForm

func (c *Context) initFormCache() {
	if c.formCache == nil {
		c.formCache = make(url.Values)
		_ = c.r.ParseMultipartForm(c.srv.MaxMultipartMemory)
		c.formCache = c.r.PostForm
	}
}

func (c *Context) PostFormArray(key string) (values []string, ok bool) {
	c.initFormCache()
	values, ok = c.formCache[key]
	return
}

func (c *Context) PostForm(key string) (string, bool) {
	if values, ok := c.PostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (c *Context) PostFormString(key string) string {
	value, _ := c.PostForm(key)
	return value
}

func (c *Context) PostFormBool(key string) bool {
	value, _ := c.PostForm(key)
	return cast.ToBool(value)
}

func (c *Context) PostFormInt(key string) int {
	value, _ := c.PostForm(key)
	return cast.ToInt(value)
}

func (c *Context) PostFormInt64(key string) int64 {
	value, _ := c.PostForm(key)
	return cast.ToInt64(value)
}

func (c *Context) PostFormFloat64(key string) float64 {
	value, _ := c.PostForm(key)
	return cast.ToFloat64(value)
}

func (c *Context) FormFileHeader(name string) (*multipart.FileHeader, error) {
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

func (c *Context) GetHeader(key string) string {
	return c.r.Header.Get(key)
}

/// Binding struct

func (c *Context) JsonBinding(d any) error {
	dc := json.NewDecoder(c.r.Body)
	if err := dc.Decode(d); err != nil {
		return err
	}
	return validate.Struct(d)
}

/// Cookie

func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.r.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite) {
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

func (c *Context) Status(code int) {
	c.w.WriteHeader(code)
}

func (c *Context) Header(key, value string) {
	if value == "" {
		c.w.Header().Del(key)
		return
	}
	c.w.Header().Set(key, value)
}

func (c *Context) ResponseWithStatus(code int, data []byte) error {
	c.Status(code)
	_, err := c.w.Write(data)
	return err
}

func (c *Context) Response(data []byte) error {
	return c.ResponseWithStatus(http.StatusOK, data)
}

func (c *Context) Json(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return c.Response(b)
}

/// context.Context

func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) Err() error {
	return c.ctx.Err()
}

func (c *Context) Value(key any) any {
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	return c.ctx.Value(key)
}
