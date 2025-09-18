package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/coze-dev/coze-studio/backend/api/model/plugin_develop"
	common "github.com/coze-dev/coze-studio/backend/api/model/plugin_develop/common"
	"github.com/coze-dev/coze-studio/backend/application/plugin"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
)

// GetPlaygroundPluginListByWanwu 参考GetPlaygroundPluginList
// @router /api/plugin_api/get_playground_plugin_list [POST]
func GetPlaygroundPluginListByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req plugin_develop.GetPlaygroundPluginListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	if req.GetSpaceID() <= 0 {
		invalidParamRequestResponse(c, "spaceID is invalid")
		return
	}
	if req.GetPage() <= 0 {
		invalidParamRequestResponse(c, "page is invalid")
		return
	}
	if req.GetSize() >= 30 {
		invalidParamRequestResponse(c, "size is invalid")
		return
	}

	// when there is only one element in the types list, and the element type is workflow, use workflow service
	// TODO Figure out when there are multiple values for types
	if len(req.GetPluginTypes()) == 1 && req.GetPluginTypes()[0] == int32(common.PluginType_WORKFLOW) {
		// 修改workflow service 适配wanwu
		resp, err := appworkflow.SVC.GetPlaygroundPluginListByWanwu(ctx, &req)
		if err != nil {
			internalServerErrorResponse(ctx, c, err)
			return
		}
		c.JSON(consts.StatusOK, resp)
		return
	}

	resp, err := plugin.PluginApplicationSVC.GetPlaygroundPluginList(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}
