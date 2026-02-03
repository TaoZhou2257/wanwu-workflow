package coze

import (
	"context"
	"strconv"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	user_entity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
	tConsts "github.com/coze-dev/coze-studio/backend/types/consts"
)

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
	//处理 parameters 字段，防止 JSON "null" 导致 map 变成 nil
	if req.Parameters == nil || *req.Parameters == "" || *req.Parameters == "null" {
		emptyJSON := "{}"
		req.Parameters = &emptyJSON
	}
	sr, err := appworkflow.SVC.OpenAPIChatFlowRun(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	sendChatFlowStreamRunSSE(ctx, w, sr)
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
