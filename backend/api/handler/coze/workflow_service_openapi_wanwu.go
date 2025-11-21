package coze

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

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
// 0. FIXME 智能体运行该接口，不会在header中带userId、orgId，需要在该方法中设置ctxcache
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

// OpenAPICreateConversationByWanwu参考 OpenAPICreateConversation
// @router /v1/workflow/conversation/create_by_wanwu [POST]
func OpenAPICreateConversationByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.CreateConversationRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	userID := ctxutil.GetApiAuthFromCtx(ctx).UserID
	newAppID, _ := appworkflow.SVC.IDGenerator.GenID(ctx)
	// 创建conversation template草稿
	_, err = appworkflow.GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
		AppID:   newAppID,
		UserID:  userID,
		SpaceID: mustParseInt64(*req.SpaceID),
		Name:    *req.ConversationMame,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	req.AppID = ptr.Of(strconv.FormatInt(newAppID, 10))
	resp, err := appworkflow.SVC.OpenAPICreateConversation(ctx, &req)
	// 将appId返回给bff
	resp.ConversationData.MetaData = map[string]string{
		"appId": strconv.FormatInt(newAppID, 10),
	}
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}

	c.JSON(consts.StatusOK, resp)
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

func mustParseInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}
