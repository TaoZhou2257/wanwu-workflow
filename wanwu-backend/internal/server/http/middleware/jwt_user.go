package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/config"
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/pkg/httputil"
	jwt_util "github.com/UnicomAI/wanwu/pkg/jwt-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func JwtUser(ctx context.Context, appCtx *app.RequestContext) {
	requestAuthType := appCtx.GetInt32(middleware.RequestAuthTypeStr)
	if requestAuthType != int32(middleware.RequestAuthTypeWebAPI) {
		httputil.InternalError(ctx, appCtx, errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "invalid request auth type")))
		return
	}

	token, err := getJWTToken(appCtx)
	if err != nil {
		httputil.InternalError(ctx, appCtx, errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", err.Error())))
		return
	}

	claims, err := jwt_util.ParseToken(token)
	if err != nil {
		httputil.InternalError(ctx, appCtx, errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", err.Error())))
		return
	}
	if claims.Subject != jwt_util.USER {
		httputil.InternalError(ctx, appCtx, errorx.New(errno.ErrUserAuthenticationFailed, errorx.KV("reason", "invalid token subject")))
		return
	}

	orgID := appCtx.Request.Header.Get(config.X_ORG_ID)
	ctxcache.Store(ctx, config.X_ORG_ID, orgID)

	ctxcache.Store(ctx, consts.SessionDataKeyInCtx, &entity.Session{
		UserID:    util.MustI64(claims.UserID),
		Locale:    string(i18n.GetLocale(ctx)),
		CreatedAt: time.Unix(claims.NotBefore, 0),
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	})
	appCtx.Next(ctx)
}

// 从Header Authorization中获取Token
func getJWTToken(ctx *app.RequestContext) (token string, err error) {
	authorization := ctx.Request.Header.Get("Authorization")
	if authorization != "" {
		tks := strings.Split(authorization, " ")
		if len(tks) > 1 && tks[0] == "Bearer" {
			return tks[1], err
		} else {
			err = fmt.Errorf("not Bearer token format")
			return "", err
		}
	} else {
		return "", fmt.Errorf("token is nil")
	}
}
