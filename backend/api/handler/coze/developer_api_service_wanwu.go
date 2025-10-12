package coze

import (
	"context"
	"net/url"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/application/upload"
)

// GetIconByWanwu 参考GetIcon 返回默认图片
// @router /api/developer/get_icon [POST]
func GetIconByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req developer_api.GetIconRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := upload.SVC.GetIcon(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	for _, icon := range resp.Data.IconList {
		// 设置默认图标
		icon.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
			os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))

	}
	c.JSON(consts.StatusOK, resp)
}
