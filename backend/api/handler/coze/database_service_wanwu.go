package coze

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/model/data/database/table"
	"github.com/coze-dev/coze-studio/backend/application/memory"
)

// ListDatabaseByWanwu参考 ListDatabase
// @router /api/memory/database/list [POST]
func ListDatabaseByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req table.ListDatabaseRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := memory.DatabaseApplicationSVC.ListDatabase(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	// Adjust icon URL for Wanwu
	baseUrl := os.Getenv("WANWU_EXTERNAL_SCHEME") + "://" + os.Getenv("WANWU_EXTERNAL_ENDPOINT")
	url, _ := url.JoinPath(baseUrl, "/api/static/icon/icon-Database-v2.jpg")
	for _, db := range resp.DatabaseInfoList {
		if strings.Contains(db.IconURL, "default_database_icon.png") {
			db.IconURL = url
		}
	}
	c.JSON(consts.StatusOK, resp)
}

// GetDatabaseByIDByWanwu 参考GetDatabaseByID
// @router /api/memory/database/get_by_id [POST]
func GetDatabaseByIDByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req table.SingleDatabaseRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := memory.DatabaseApplicationSVC.GetDatabaseByID(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	// Adjust icon URL for Wanwu
	baseUrl := os.Getenv("WANWU_EXTERNAL_SCHEME") + "://" + os.Getenv("WANWU_EXTERNAL_ENDPOINT")
	url, _ := url.JoinPath(baseUrl, "/api/static/icon/icon-Database-v2.jpg")
	if strings.Contains(resp.DatabaseInfo.IconURL, "default_database_icon.png") {
		resp.DatabaseInfo.IconURL = url
	}
	c.JSON(consts.StatusOK, resp)
}

// GetConnectorNameByWanwu 参考GetConnectorName
// @router /api/memory/database/get_connector_name [POST]
func GetConnectorNameByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req table.GetSpaceConnectorListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := memory.DatabaseApplicationSVC.GetConnectorNameByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}
