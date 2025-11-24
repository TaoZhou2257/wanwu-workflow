package coze

import (
	"context"
	"encoding/base64"
	"net/url"
	"os"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/application/upload"
	crossupload "github.com/coze-dev/coze-studio/backend/crossdomain/upload"
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

// UploadFileByWanwu 参考UploadFile 通过文件名、base64格式上传文件
// @router /api/bot/upload_file_by_wanwu [POST]
func UploadFileByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req uploadFileRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	fileContent, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	uploadResp, err := crossupload.DefaultWanwuSVC().UploadFileByByte(ctx, req.Name, fileContent)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, &crossupload.UploadFileResp{URL: uploadResp.URL, URI: uploadResp.URI})
}

type uploadFileRequest struct {
	// file name
	Name string `thrift:"name,1" form:"name" json:"name" query:"name"`
	// file data
	Data string `thrift:"data,2" form:"data" json:"data" query:"data"`
}
