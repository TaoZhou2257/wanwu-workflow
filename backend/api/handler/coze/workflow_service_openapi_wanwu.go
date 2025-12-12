package coze

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
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
// 3. ctx中设置WANWU_WORKFLOW_OPENAPI_RUN_RECORD_EXECUTE_HISTORY
// 4. 返回resp.Data unmarshal的结构体
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

	if os.Getenv("WANWU_WORKFLOW_OPENAPI_RUN_SKIP_EXECUTE_HISTORY") == "1" {
		ctx = context.WithValue(ctx, "WANWU_WORKFLOW_OPENAPI_RUN_SKIP_EXECUTE_HISTORY", true)
	}

	resp, tPlan, err := appworkflow.SVC.OpenAPIRunByWanwu(ctx, workflowID, &req)
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
	if tPlan == vo.ReturnVariables {
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
	} else {
		c.JSON(consts.StatusOK, resp.Data)
		return
	}

	internalServerErrorResponse(ctx, c, errors.New("empty response"))
}

// OpenAPICreateConversationByWanwu 参考OpenAPICreateConversation
// @router /v1/workflow/conversation/create_by_wanwu [POST]
func OpenAPICreateConversationByWanwu(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.CreateConversationRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	// appID非空，用于对话流调试(appID == workflowID)
	// appID为空，自动生成一个新的，用于应用广场新建对话流
	if req.AppID == nil {
		newAppID, _ := appworkflow.SVC.IDGenerator.GenID(ctx)
		req.AppID = ptr.Of(strconv.FormatInt(newAppID, 10))
	}
	resp, err := appworkflow.SVC.OpenAPICreateConversationByWanwu(ctx, &req)
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
