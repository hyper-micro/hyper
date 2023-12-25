package gin

import (
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
