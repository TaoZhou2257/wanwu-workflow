package coze

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/data/variable/project_memory"
)

// GetMemoryVariableMetaByWanwu 参考GetMemoryVariableMeta
// @router /api/memory/variable/get_meta [POST]
func GetMemoryVariableMetaByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req project_memory.GetMemoryVariableMetaReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := &project_memory.GetMemoryVariableMetaResp{
		VariableMap: nil,
		Code:        0,
		Msg:         "",
		BaseResp:    nil,
	}
	c.JSON(consts.StatusOK, resp)
}

// GetProjectVariableListByWanwu 参考GetProjectVariableList
// @router /api/memory/project/variable/meta_list [GET]
func GetProjectVariableListByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req project_memory.GetProjectVariableListReq
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp := &project_memory.GetProjectVariableListResp{
		VariableList: nil,
		CanEdit:      false,
		GroupConf:    nil,
		Code:         0,
		Msg:          "",
		BaseResp:     nil,
	}

	c.JSON(consts.StatusOK, resp)
}
