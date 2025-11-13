package middleware

import (
	"context"

	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/config"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
)

func SetOrgID(ctx context.Context, appCtx *app.RequestContext) {
	// orgID
	orgID := appCtx.Request.Header.Get(config.X_ORG_ID)
	if orgID != "" {
		ctxcache.Store(ctx, config.X_ORG_ID, orgID)
	}
	appCtx.Next(ctx)
}
