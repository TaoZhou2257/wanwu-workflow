package middleware

import (
	"context"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

// ！！！同步于：backend/api/middleware/host.go
// ！！！同步于：backend/api/middleware/host.go
// ！！！同步于：backend/api/middleware/host.go

// SetHost 参考SetHostMW
func SetHost(ctx context.Context, appCtx *app.RequestContext) {
	ctxcache.Store(ctx, consts.HostKeyInCtx, os.Getenv("WANWU_EXTERNAL_ENDPOINT"))
	ctxcache.Store(ctx, consts.RequestSchemeKeyInCtx, string(appCtx.GetRequest().Scheme()))
	appCtx.Next(ctx)
}
