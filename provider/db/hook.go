package db

import (
	"context"

	"xorm.io/xorm/contexts"
)

type Hook struct {
	host     string
	port     int
	database string
}

func (Hook) BeforeProcess(c *contexts.ContextHook) (context.Context, error) {
	return c.Ctx, nil
}

func (h Hook) AfterProcess(c *contexts.ContextHook) error {
	//sp, _ := sdk.CSpan(c.Ctx, nil, nil, porter.TagMySQLClient).Start()
	//sp.Describe(porter.Describe{
	//	ToSvcName: porter.String(h.host),
	//	ToSvcPort: porter.Int64(cast.ToInt64(h.port)),
	//	ToMethod:  porter.String(h.database),
	//})
	//sp.Metadata(porter.MD{
	//	"isSlow": c.ExecuteTime > time.Second,
	//})
	//
	//sp.Fields(porter.MD{
	//	"sql":         c.SQL,
	//	"executeTime": c.ExecuteTime.Milliseconds(),
	//	"args":        c.Args,
	//})
	//sp.MaybeError(c.Err).Sync()
	return nil
}
