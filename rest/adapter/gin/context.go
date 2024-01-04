package gin

import (
	"mime/multipart"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyper-micro/hyper/rest/router"
)

type ginAdapterContext struct {
	ginCtx *gin.Context
}

func newGinAdapterContext(ginCtx *gin.Context) router.Context {
	return &ginAdapterContext{
		ginCtx: ginCtx,
	}
}

func (c *ginAdapterContext) Deadline() (deadline time.Time, ok bool) {
	return c.ginCtx.Deadline()
}

func (c *ginAdapterContext) Done() <-chan struct{} {
	return c.ginCtx.Done()
}

func (c *ginAdapterContext) Err() error {
	return c.ginCtx.Err()
}

func (c *ginAdapterContext) Value(key any) any {
	return c.ginCtx.Value(key)
}

func (c *ginAdapterContext) Request() *http.Request {
	return c.ginCtx.Request
}

func (c *ginAdapterContext) Writer() http.ResponseWriter {
	return c.ginCtx.Writer
}

func (c *ginAdapterContext) GetRawData() ([]byte, error) {
	return c.ginCtx.GetRawData()
}

func (c *ginAdapterContext) Next() {
	c.ginCtx.Next()
}

func (c *ginAdapterContext) Abort() {
	c.ginCtx.Abort()
}

func (c *ginAdapterContext) FullPath() string {
	return c.ginCtx.FullPath()
}

func (c *ginAdapterContext) Set(key string, value any) {
	c.ginCtx.Set(key, value)
}

func (c *ginAdapterContext) Get(key string) (any, bool) {
	return c.ginCtx.Get(key)
}

func (c *ginAdapterContext) Param(key string) string {
	return c.ginCtx.Param(key)
}

func (c *ginAdapterContext) Query(key string) string {
	return c.ginCtx.Query(key)
}

func (c *ginAdapterContext) PostForm(key string) string {
	return c.ginCtx.PostForm(key)
}

func (c *ginAdapterContext) PostFormArray(key string) []string {
	return c.ginCtx.PostFormArray(key)
}

func (c *ginAdapterContext) FormFile(name string) (*multipart.FileHeader, error) {
	return c.ginCtx.FormFile(name)
}

func (c *ginAdapterContext) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool, sameSite http.SameSite) {
	c.ginCtx.SetSameSite(sameSite)
	c.ginCtx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func (c *ginAdapterContext) Cookie(name string) (string, error) {
	return c.ginCtx.Cookie(name)
}

func (c *ginAdapterContext) Status(code int) {
	c.ginCtx.Status(code)
}

func (c *ginAdapterContext) Header(key, value string) {
	c.ginCtx.Header(key, value)
}

func (c *ginAdapterContext) GetHeader(key string) string {
	return c.ginCtx.GetHeader(key)
}

func (c *ginAdapterContext) Response(code int, contentType string, data []byte) {
	c.ginCtx.Data(code, contentType, data)
}
