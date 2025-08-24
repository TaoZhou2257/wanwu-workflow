package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	upload "github.com/coze-dev/coze-studio/backend/api/model/file/upload"
	uploadSVC "github.com/coze-dev/coze-studio/backend/application/upload"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	coze_consts "github.com/coze-dev/coze-studio/backend/types/consts"
)

// ApplyUploadActionByWanwu 参考ApplyUploadAction，修改host，workflow-wanwu:8999 => <external_ip>:8081
// @router /api/common/upload/apply_upload_action [GET]
func ApplyUploadActionByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req upload.ApplyUploadActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := new(upload.ApplyUploadActionResponse)
	host, ok := ctxcache.Get[string](ctx, coze_consts.HostKeyInCtx)
	if !ok {
		host = string(c.Request.Host())
	}
	if ptr.From(req.Action) == "ApplyImageUpload" {
		resp, err = uploadSVC.SVC.ApplyImageUpload(ctx, &req, host)
		if err != nil {
			internalServerErrorResponse(ctx, c, err)
			return
		}
	} else if ptr.From(req.Action) == "CommitImageUpload" {
		resp, err = uploadSVC.SVC.CommitImageUpload(ctx, &req, host)
		if err != nil {
			internalServerErrorResponse(ctx, c, err)
			return
		}
	}

	c.JSON(consts.StatusOK, resp)
}
