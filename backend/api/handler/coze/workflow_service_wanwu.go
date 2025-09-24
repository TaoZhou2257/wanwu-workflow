package coze

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/upload"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

// GetWorkFlowListByWanwu 参考GetWorkFlowList
// @router /api/workflow_api/workflow_list_by_wanwu
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
// @router /api/workflow_api/workflow_select_by_wanwu
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

// GetWorkFlowOpenAPIV3SchemaByWanwu 获取workflow openapi v3 schema
// @router /v1/workflow/:workflow_id/schema_by_wanwu [GET]
func GetWorkFlowOpenAPIV3SchemaByWanwu(ctx context.Context, c *app.RequestContext) {
	workflowID := c.Param("workflow_id")
	if workflowID == "" {
		invalidParamRequestResponse(c, "workflow_id empty")
		return
	}
	wfSchema, err := appworkflow.SVC.GetWorkFlowOpenAPIV3SchemaByWanwu(ctx, workflowID)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	c.JSON(consts.StatusOK, wfSchema)
}

// OpenAPIRunWorkFlowByWanwu 参考OpenAPIRunFlow
// 0. FIXME 智能体运行该接口，不会在header中带userId、orgId，跳过jwt校验后，需要在该方法中设置ctxcache
// 1. 将workflow_id从 body => path
// 2. 将body参数{...} marsharl到req.Parameters上
// 3. 返回resp.Data unmarshal的结构体
// @router /v1/workflow/:workflow_id/run_by_wanwu [POST]
func OpenAPIRunWorkFlowByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error

	workflowID := c.Param("workflow_id")
	if workflowID == "" {
		invalidParamRequestResponse(c, "workflow_id empty")
		return
	}

	parameters, err := preprocessWorkflowRequestBodyByWanwu(ctx, c)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	} else {

	}

	var req workflow.OpenAPIRunFlowRequest

	if err := c.BindAndValidate(&req); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	req.Parameters = parameters

	resp, err := appworkflow.SVC.OpenAPIRunByWanwu(ctx, workflowID, &req)
	if err != nil {
		var se vo.WorkflowError
		if errors.As(err, &se) {
			resp = new(workflow.OpenAPIRunFlowResponse)
			resp.Code = int64(se.OpenAPICode())
			resp.Msg = ptr.Of(se.Msg())
			debugURL := se.DebugURL()
			if debugURL != "" {
				resp.DebugUrl = ptr.Of(debugURL)
			}
			c.JSON(consts.StatusOK, resp)
			return
		}

		internalServerErrorResponse(ctx, c, err)
		return
	}

	var respData map[string]any
	if resp.Data != nil {
		if err = sonic.Unmarshal([]byte(*resp.Data), &respData); err != nil {
			logs.CtxErrorf(ctx, "unmarshal resp.Data (%v) err: %v", resp.Data, err)
			c.JSON(consts.StatusOK, resp.Data)
			return
		}
		c.JSON(consts.StatusOK, respData)
		return
	}
	internalServerErrorResponse(ctx, c, errors.New("empty response"))
}

// preprocessWorkflowRequestBodyByWanwu 参考preprocessWorkflowRequestBody
func preprocessWorkflowRequestBodyByWanwu(_ context.Context, c *app.RequestContext) (*string, error) {
	// Read the raw request body
	rawData, err := c.Request.BodyE()
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	// Unmarshal into a temporary map
	var bodyData map[string]interface{}
	if err = sonic.Unmarshal(rawData, &bodyData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal request body: %w", err)
	}

	return ptr.Of(string(rawData)), nil
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

// GetIconByWanwu . 参考GetIcon 返回默认图片
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

// ImportWorkFlow .
// @router /api/workflow_api/import [POST]
func ImportWorkFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req importWorkflowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	// create
	createReq := workflow.CreateWorkflowRequest{
		SpaceID: req.SpaceID,
		Name:    req.Name,
		Desc:    req.Desc,
		IconURI: "default_icon/default_workflow_icon.png",
	}
	resp, err := appworkflow.SVC.CreateWorkflow(ctx, &createReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	// save
	saveReq := workflow.SaveWorkflowRequest{
		WorkflowID: resp.Data.WorkflowID,
		SpaceID:    &req.SpaceID,
		Schema:     &req.Schema,
	}
	_, err = appworkflow.SVC.SaveWorkflow(ctx, &saveReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, resp)
}

type importWorkflowRequest struct {
	// Space id, cannot be empty
	SpaceID string `form:"space_id,required" json:"space_id,required" query:"space_id,required"`
	// process name
	Name string `form:"name,required" json:"name,required" query:"name,required"`
	// Process description, not null
	Desc string `form:"desc,required" json:"desc,required" query:"desc,required"`
	// file data
	Schema string `form:"schema" json:"schema" query:"schema"`
}
