package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
)

func I18n(ctx context.Context, appCtx *app.RequestContext) {
	acceptLanguage := string(appCtx.Request.Header.Get("Accept-Language"))
	locale := string(i18n.LocaleZH)
	if acceptLanguage != "" {
		languages := strings.Split(acceptLanguage, ",")
		if len(languages) > 0 {
			locale = languages[0]
		}
	}

	ctx = i18n.SetLocale(ctx, locale)

	appCtx.Next(ctx)
}
