package coze

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	user_entity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
	tConsts "github.com/coze-dev/coze-studio/backend/types/consts"
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

// UpdateWorkflowMetaByWanwu 参考UpdateWorkflowMeta
// @router /api/workflow_api/update_meta_by_wanwu [POST]
func UpdateWorkflowMetaByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.UpdateWorkflowMetaRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.UpdateWorkflowMetaByWanwu(ctx, &req)

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
		workflowDefaultIconURL(workflowData)
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
		workflowDefaultIconURL(workflowData)
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
	for _, workflowDetail := range workflowDetailInfoDataList.List {
		workflowDetailDefaultIconURL(workflowDetail)
	}

	response := map[string]any{
		"data":    workflowDetailInfoDataList,
		"code":    0,
		"message": "",
	}

	c.JSON(consts.StatusOK, response)
}

// GetCanvasInfoByWanwu 参考GetCanvasInfo 增加返回图标判断
// @router /api/workflow_api/canvas/draft [POST]
func GetCanvasInfoByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req appworkflow.ExportWorkflowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	req.Type = workflowModel.FromDraft
	resp, err := appworkflow.SVC.GetCanvasInfoByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	workflowDefaultIconURL(resp.Data.Workflow)

	c.JSON(consts.StatusOK, resp)
}

// OpenAPIGetWorkflowInfoByWanwu 参考OpenAPIGetWorkflowInfo
// 0. FIXME 前端运行该接口，不会在header中带orgId，需要在该方法中设置ctxcache
// @router /v1/workflows/:workflow_id [GET]
func OpenAPIGetWorkflowInfoByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error

	if err = processOpenAPIGetWorkflowInfoRequest(ctx, c); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	var req workflow.OpenAPIGetWorkflowInfoRequest

	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置ctxcache
	wf, err := appworkflow.GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{
		ID:       mustParseInt64(req.GetWorkflowID()),
		MetaOnly: true,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	if _, ok := ctxcache.Get[string](ctx, "X-Org-Id"); !ok {
		ctxcache.Store(ctx, "X-Org-Id", strconv.Itoa(int(wf.Meta.SpaceID)))
	}

	resp, err := appworkflow.SVC.OpenAPIGetWorkflowInfo(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// OpenAPIChatFlowRunByWanwu 参考OpenAPIChatFlowRun
// 0. FIXME 前端运行该接口，不会在header中带userId、orgId，需要在该方法中设置ctxcache
// @router /v1/workflows/chat [POST]
func OpenAPIChatFlowRunByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	if err = preprocessWorkflowRequestBody(ctx, c); err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	var req workflow.ChatFlowRunRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	// 设置ctxcache
	if ctxutil.GetUserSessionFromCtx(ctx) == nil {
		meta, err := appworkflow.GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{
			ID:       mustParseInt64(req.WorkflowID),
			MetaOnly: true,
		})
		if err != nil {
			internalServerErrorResponse(ctx, c, err)
			return
		}
		ctxcache.Store(ctx, tConsts.SessionDataKeyInCtx, &user_entity.Session{
			UserID: meta.CreatorID,
			Locale: string(i18n.GetLocale(ctx)),
		})
		ctxcache.Store(ctx, "X-Org-Id", strconv.Itoa(int(meta.SpaceID)))
	}

	w := sse.NewWriter(c)
	c.SetContentType("text/event-stream; charset=utf-8")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")
	c.Response.Header.Set("Access-Control-Allow-Origin", "*")

	sr, err := appworkflow.SVC.OpenAPIChatFlowRun(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	sendChatFlowStreamRunSSE(ctx, w, sr)

}

// GetWorkflowVersionListByWanwu 获取工作流版本列表
// @router /api/workflow_api/version_list [POST]
func GetWorkflowVersionListByWanwu(ctx context.Context, c *app.RequestContext) {
	var req appworkflow.GetWorkflowRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.GetWorkflowVersionListByWanwu(ctx, req.WorkflowID)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, resp)
}

// UpdateWorkflowVersionDescriptionByWanwu 更新版本描述
// @router /api/workflow_api/version/description [PUT]
func UpdateWorkflowVersionDescriptionByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req appworkflow.UpdateWorkflowVersionDescriptionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := appworkflow.SVC.UpdateWorkflowVersionDescriptionByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, resp)
}

// RollbackWorkflowVersionByWanwu 回滚工作流版本
// @router /api/workflow_api/revert [POST]
func RollbackWorkflowVersionByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req appworkflow.RollbackWorkflowVersionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.RollbackWorkflowVersionByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetWorkflowLatestVersionCanvasInfoByWanwu 获取已发布工作流最新版本信息
// @router /api/workflow_api/canvas [POST]
func GetWorkflowLatestVersionCanvasInfoByWanwu(ctx context.Context, c *app.RequestContext) {
	var req appworkflow.ExportWorkflowRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.GetCanvasInfoByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	workflowDefaultIconURL(resp.Data.Workflow)
	c.JSON(consts.StatusOK, resp)
}

// RunWorkFlowLatestVersionByWanwu 运行已发布最新版工作流
// @router /v1/workflow/run_by_wanwu [POST]
func RunWorkFlowLatestVersionByWanwu(ctx context.Context, c *app.RequestContext) {

	var err error
	var req workflow.WorkFlowTestRunRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := appworkflow.SVC.LatestVersionRunByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// GetHistorySchemaByWanwu 获取历史工作流schema .
// @router /api/workflow_api/history_schema [POST]
func GetHistorySchemaByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetHistorySchemaRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := appworkflow.SVC.GetWorkflowVersionSchema(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// CreateProjectConversationDefByWanwu 参考CreateProjectConversationDef
// @router /api/workflow_api/project_conversation/create_by_wanwu [POST]
func CreateProjectConversationDefByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.CreateProjectConversationDefRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	resp, err := appworkflow.SVC.CreateApplicationConversationDefByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// --- internal ---

func workflowDefaultIconURL(wf *workflow.Workflow) {
	if wf == nil {
		return
	}
	if wf.URL == "" || strings.Contains(wf.URL, "default_workflow_icon.png") {
		switch wf.FlowMode {
		case workflow.WorkflowMode_Workflow:
			wf.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"), os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
		case workflow.WorkflowMode_ChatFlow:
			wf.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"), os.Getenv("WANWU_CHATFLOW_DEFAULT_ICON"))
		}
	}
}

func workflowDetailDefaultIconURL(wf *workflow.WorkflowDetailInfoData) {
	if wf == nil {
		return
	}
	if wf.Icon == "" || strings.Contains(wf.Icon, "default_workflow_icon.png") {
		switch wf.FlowMode {
		case workflow.WorkflowMode_Workflow:
			wf.Icon, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"), os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
		case workflow.WorkflowMode_ChatFlow:
			wf.Icon, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"), os.Getenv("WANWU_CHATFLOW_DEFAULT_ICON"))
		}
	}
}
