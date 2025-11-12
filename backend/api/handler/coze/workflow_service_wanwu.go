package coze

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/app/intelligence"
	"github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/common"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
)

// CreateWorkflow .
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
		switch workflowData.FlowMode {
		case workflow.WorkflowMode_Workflow:
			if workflowData.URL == "" || strings.Contains(workflowData.URL, "default_workflow_icon.png") {
				// 设置默认图标
				workflowData.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
					os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
			}
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

// GetDraftIntelligenceListByWanwu 参考GetDraftIntelligenceList
// @router /api/intelligence_api/search/get_draft_intelligence_list [POST]
func GetDraftIntelligenceListByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req intelligence.GetDraftIntelligenceListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}

	headerReferer := c.Request.Header.Get("Referer")
	if headerReferer == "" {
		invalidParamRequestResponse(c, "empty header referer")
		return
	}
	parsedUrl, err := url.Parse(headerReferer)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	parsedParams, err := url.ParseQuery(parsedUrl.RawQuery)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	workflowID, _ := strconv.Atoi(parsedParams.Get("workflow_id"))

	var intellData []*intelligence.IntelligenceData
	intellData = append(intellData, &intelligence.IntelligenceData{
		BasicInfo: &common.IntelligenceBasicInfo{
			ID:          int64(workflowID),
			Name:        "todo",
			Description: "todo",
			IconURI:     "default_icon/default_agent_icon.png",
			IconURL:     "http://localhost:8888/local_storage/opencoze/default_icon/default_agent_icon.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T024057Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=693456cdb3b8121854b40670a9d5c9d5b1c82996329a2e7afdd9ddefcabc9905",
			SpaceID:     7567337229379436544,
			OwnerID:     7567337229371047936,
			CreateTime:  1762248076,
			UpdateTime:  1762254567,
			Status:      1,
			PublishTime: 0,
		},
		Type: 2,
		PublishInfo: &intelligence.IntelligencePublishInfo{
			PublishTime:  "",
			HasPublished: false,
			Connectors:   nil,
		},
		PermissionInfo: &intelligence.IntelligencePermissionInfo{
			InCollaboration: false,
			CanDelete:       true,
			CanView:         true,
		},
		OwnerInfo: &common.User{
			UserID:         7567337229371047936,
			Nickname:       "",
			AvatarURL:      "",
			UserUniqueName: "",
			UserLabel:      nil,
		},
		LatestAuditInfo: nil,
		FavoriteInfo: &intelligence.FavoriteInfo{
			IsFav:   false,
			FavTime: "",
		},
		OtherInfo: nil,
	})
	resp := &intelligence.GetDraftIntelligenceListResponse{
		Data: &intelligence.DraftIntelligenceListData{
			Intelligences: intellData,
			Total:         1,
			HasMore:       false,
			NextCursorID:  "",
		},
	}
	c.JSON(consts.StatusOK, resp)
}

// GetDraftIntelligenceInfoByWanwu 参考GetDraftIntelligenceInfo
// @router /api/intelligence_api/search/get_draft_intelligence_info [POST]
func GetDraftIntelligenceInfoByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req intelligence.GetDraftIntelligenceInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp := &intelligence.GetDraftIntelligenceInfoResponse{
		Data: &intelligence.GetDraftIntelligenceInfoData{
			IntelligenceType: 2,
			BasicInfo: &common.IntelligenceBasicInfo{
				ID:          req.IntelligenceID,
				Name:        "todo",
				Description: "todo",
				IconURI:     "defaul t_icon/defaul t_agent_icon. png",
				IconURL:     "http://localhost:8888/local_storage/opencoze/default_icon/default_agent_icon.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20251105%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20251105T024057Z&X-Amz-Expires=604800&X-Amz-SignedHeaders=host&X-Amz-Signature=693456cdb3b8121854b40670a9d5c9d5b1c82996329a2e7afdd9ddefcabc9905",
				SpaceID:     7567337229379436544,
				OwnerID:     7567337229371047936,
				CreateTime:  1762248076,
				UpdateTime:  0,
				Status:      1,
				PublishTime: 1762248076,
			},
			PublishInfo: &intelligence.IntelligencePublishInfo{
				PublishTime:  "",
				HasPublished: false,
				Connectors:   nil,
			},
			OwnerInfo: &common.User{
				UserID:         7567337229371047936,
				Nickname:       "",
				AvatarURL:      "",
				UserUniqueName: "",
				UserLabel:      nil,
			},
		},
		Code:     0,
		Msg:      "",
		BaseResp: nil,
	}
	c.JSON(consts.StatusOK, resp)
}

// ListProjectConversationDefByWanwu 参考 ListProjectConversationDef
// @router /api/workflow_api/project_conversation/list [GET]
func ListProjectConversationDefByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.ListProjectConversationRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp := &workflow.ListProjectConversationResponse{
		Data:     nil,
		Cursor:   "",
		Code:     0,
		Msg:      "",
		BaseResp: nil,
	}
	if req.GetCreateMethod() == workflow.CreateMethod_ManualCreate {
		var wpc []*workflow.ProjectConversation
		wpc = append(wpc, &workflow.ProjectConversation{
			UniqueID:                "",
			ConversationName:        "Default",
			ConversationID:          strconv.FormatInt(time.Now().UnixMilli(), 10),
			ReleaseConversationName: "",
		})
		resp.Data = wpc
	}
	c.JSON(consts.StatusOK, resp)
}

// OpenAPICreateConversationByWanwu 参考OpenAPICreateConversationByWanwu
// @router /v1/workflow/conversation/create [POST]
func OpenAPICreateConversationByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.CreateConversationRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := appworkflow.SVC.OpenAPICreateConversationByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}

// OpenAPIGetWorkflowInfoByWanwu 参考OpenAPIGetWorkflowInfoByWanwu
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

	resp, err := appworkflow.SVC.OpenAPIGetWorkflowInfo(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
}
