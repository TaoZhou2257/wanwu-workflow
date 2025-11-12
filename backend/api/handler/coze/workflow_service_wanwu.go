package coze

import (
	"context"
	"net/url"
	"os"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
)

// CreateWorkflowByWanwu 参考CreateWorkflow
// @router /api/workflow_api/create [POST]
func CreateWorkflowByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.CreateWorkflowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := appworkflow.SVC.CreateWorkflowByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// CopyWorkflowByWanwu 参考CopyWorkflow
// @router /api/workflow_api/copy [POST]
func CopyWorkflowByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.CopyWorkflowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.CopyWorkflowByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetWorkFlowListByWanwu 参考GetWorkFlowList
// @router /api/workflow_api/workflow_list_by_wanwu [POST]
func GetWorkFlowListByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetWorkFlowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.ListWorkflowByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	for _, workflowData := range resp.Data.WorkflowList {
		if workflowData.URL == "" || strings.Contains(workflowData.URL, "default_workflow_icon.png") {
			// 设置默认图标
			workflowData.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
				os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
		}
	}
	c.JSON(consts.StatusOK, resp)
}

// GetWorkFlowSelectByWanwu 参考GetWorkFlowList
// @router /api/workflow_api/workflow_select_by_wanwu [POST]
func GetWorkFlowSelectByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetWorkFlowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.GetWorkFlowSelectByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	for _, workflowData := range resp.Data.WorkflowList {
		if workflowData.URL == "" || strings.Contains(workflowData.URL, "default_workflow_icon.png") {
			// 设置默认图标
			workflowData.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
				os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
		}
	}

	c.JSON(consts.StatusOK, resp)
}

// GetExampleWorkFlowListByWanwu 参考GetExampleWorkFlowList，目前去掉其中的example
// @router /api/workflow_api/example_workflow_list [POST]
func GetExampleWorkFlowListByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetExampleWorkFlowListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp := &workflow.GetExampleWorkFlowListResponse{
		Data: &workflow.WorkFlowListData{
			AuthList:     make([]*workflow.ResourceAuthInfo, 0),
			WorkflowList: make([]*workflow.Workflow, 0),
		},
	}

	c.JSON(consts.StatusOK, resp)
}

// ListWorkFlowOpenAPIV3SchemaByWanwu 获取workflow list openapi v3 schema
// @router /v1/workflow/list_schema_by_wanwu [POST]
func ListWorkFlowOpenAPIV3SchemaByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetWorkflowDetailRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	schemas, err := appworkflow.SVC.ListWorkFlowOpenAPIV3SchemaByWanwu(ctx, req.WorkflowIds)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	c.JSON(consts.StatusOK, schemas)
}

// GetWorkflowDetailInfoByWanwu 参考GetWorkflowDetailInfo 替换返回URL
// @router /api/workflow_api/workflow_detail_info [POST]
func GetWorkflowDetailInfoByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetWorkflowDetailInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	workflowDetailInfoDataList, err := appworkflow.SVC.GetWorkflowDetailInfo(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	for _, workflowData := range workflowDetailInfoDataList.List {
		if workflowData.Icon == "" || strings.Contains(workflowData.Icon, "default_workflow_icon.png") {
			// 设置默认图标
			workflowData.Icon, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
				os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
		}
	}

	response := map[string]any{
		"data":    workflowDetailInfoDataList,
		"code":    0,
		"message": "",
	}

	c.JSON(consts.StatusOK, response)
}

// GetCanvasInfoByWanwu 参考GetCanvasInfo 增加返回图标判断
// @router /api/workflow_api/canvas [POST]
func GetCanvasInfoByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetCanvasInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := appworkflow.SVC.GetCanvasInfo(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	if resp.Data.Workflow.URL == "" || strings.Contains(resp.Data.Workflow.URL, "default_workflow_icon.png") {
		// 设置默认图标
		resp.Data.Workflow.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
			os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
	}

	c.JSON(consts.StatusOK, resp)
}
