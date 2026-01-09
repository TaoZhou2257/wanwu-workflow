package coze

import (
	"context"
	"net/url"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	dataset "github.com/coze-dev/coze-studio/backend/api/model/data/knowledge"
)

// GetIconForDatasetByWanwu 参考GetIconForDataset
// @router /api/knowledge/icon/get [POST]
func GetIconForDatasetByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req dataset.GetIconRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	baseUrl := os.Getenv("WANWU_EXTERNAL_SCHEME") + "://" + os.Getenv("WANWU_EXTERNAL_ENDPOINT")

	resp := new(dataset.GetIconResponse)
	url, _ := url.JoinPath(baseUrl, "/api/static/icon/icon-Database-v2.jpg")
	resp.Icon = &dataset.Icon{
		URL: url,
		URI: "default_icon/default_database_icon.png",
	}
	c.JSON(consts.StatusOK, resp)
}
