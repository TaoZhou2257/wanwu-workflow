package workflow

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/api/model/base"
	"github.com/coze-dev/coze-studio/backend/api/model/playground"
	pluginAPI "github.com/coze-dev/coze-studio/backend/api/model/plugin_develop"
	"github.com/coze-dev/coze-studio/backend/api/model/plugin_develop/common"
	resource "github.com/coze-dev/coze-studio/backend/api/model/resource/common"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/application/user"
	"github.com/coze-dev/coze-studio/backend/bizpkg/debugutil"
	model "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	search "github.com/coze-dev/coze-studio/backend/domain/search/entity"
	user_entity "github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/i18n"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/maps"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/slices"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ternary"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/safego"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
	"github.com/coze-dev/coze-studio/backend/types/consts"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-resty/resty/v2"
	xmaps "golang.org/x/exp/maps"
)

// CreateWorkflowByWanwu 参考CreateWorkflow，适配了chatflow
// 1. 调换顺序，先创建workflow，再创建conversation
// 2. conversation名，默认为工作流(chatflow)名
// 3. conversation的app id，为创建的工作流id，这样conversation对应唯一的工作流(chatflow)
func (w *ApplicationService) CreateWorkflowByWanwu(ctx context.Context, req *workflow.CreateWorkflowRequest) (
	_ *workflow.CreateWorkflowResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	uID := ctxutil.MustGetUIDFromCtx(ctx)
	spaceID := mustParseInt64(req.GetSpaceID())
	if err := checkUserSpace(ctx, uID, spaceID); err != nil {
		return nil, err
	}

	wf := &vo.MetaCreate{
		CreatorID:        uID,
		SpaceID:          spaceID,
		ContentType:      workflow.WorkFlowType_User,
		Name:             req.Name,
		Desc:             req.Desc,
		IconURI:          req.IconURI,
		AppID:            parseInt64(req.ProjectID),
		Mode:             ternary.IFElse(req.IsSetFlowMode(), req.GetFlowMode(), workflow.WorkflowMode_Workflow),
		InitCanvasSchema: vo.GetDefaultInitCanvasJsonSchema(i18n.GetLocale(ctx)),
	}
	if req.IsSetFlowMode() && req.GetFlowMode() == workflow.WorkflowMode_ChatFlow {
		wf.InitCanvasSchema = vo.GetDefaultInitCanvasJsonSchemaChat(i18n.GetLocale(ctx), req.Name)
	}

	id, err := GetWorkflowDomainSVC().Create(ctx, wf)
	if err != nil {
		return nil, err
	}

	err = PublishWorkflowResource(ctx, id, ptr.Of(int32(wf.Mode)), search.Created, &search.ResourceDocument{
		Name:          &wf.Name,
		APPID:         wf.AppID,
		SpaceID:       &wf.SpaceID,
		OwnerID:       &wf.CreatorID,
		PublishStatus: ptr.Of(resource.PublishStatus_UnPublished),
		CreateTimeMS:  ptr.Of(time.Now().UnixMilli()),
	})
	if err != nil {
		return nil, vo.WrapError(errno.ErrNotifyWorkflowResourceChangeErr, err)
	}

	if req.IsSetFlowMode() && req.GetFlowMode() == workflow.WorkflowMode_ChatFlow {
		_, err := GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
			AppID:   id,
			UserID:  uID,
			SpaceID: spaceID,
			Name:    req.Name,
		})
		if err != nil {
			return nil, err
		}
		// 如果是对话流还需要创建chat_flow_role
		_, err = GetWorkflowDomainSVC().CreateChatFlowRole(ctx, &vo.ChatFlowRoleCreate{
			WorkflowID: id,
			CreatorID:  uID,
		})
		if err != nil {
			return nil, err
		}
	}

	return &workflow.CreateWorkflowResponse{
		Data: &workflow.CreateWorkflowData{
			WorkflowID: strconv.FormatInt(id, 10),
		},
	}, nil
}

// UpdateWorkflowMetaByWanwu 参考UpdateWorkflowMeta
// 1. 如果是转换成chatflow模式，且没有对话模板，则创建一个默认的对话模板
func (w *ApplicationService) UpdateWorkflowMetaByWanwu(ctx context.Context, req *workflow.UpdateWorkflowMetaRequest) (
	_ *workflow.UpdateWorkflowMetaResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	if err := checkUserSpace(ctx, ctxutil.MustGetUIDFromCtx(ctx), mustParseInt64(req.GetSpaceID())); err != nil {
		return nil, err
	}

	workflowID := mustParseInt64(req.GetWorkflowID())
	if req.IsSetFlowMode() && req.GetFlowMode() == workflow.WorkflowMode_ChatFlow {
		def, err := w.ListApplicationConversationDef(ctx, &workflow.ListProjectConversationRequest{
			ProjectID:    strconv.FormatInt(workflowID, 10),
			CreateMethod: workflow.CreateMethod_ManualCreate,
			CreateEnv:    workflow.CreateEnv_Draft,
			Limit:        1000,
			SpaceID:      req.GetSpaceID(),
		})
		if err != nil {
			return nil, err
		}
		if len(def.Data) == 0 {
			wf, err := GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{
				ID:       mustParseInt64(req.GetWorkflowID()),
				MetaOnly: true,
			})
			if err != nil {
				return nil, err
			}
			_, err = GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
				AppID:   workflowID,
				UserID:  ctxutil.MustGetUIDFromCtx(ctx),
				SpaceID: mustParseInt64(req.GetSpaceID()),
				Name:    wf.Name,
			})
			if err != nil {
				return nil, err
			}
		}
		oldRole, err := GetWorkflowDomainSVC().GetChatFlowRole(ctx, workflowID, "")
		if err != nil {
			return nil, err
		}
		if oldRole == nil {
			// 如果是对话流并且role不存在还需要创建chat_flow_role
			_, err = GetWorkflowDomainSVC().CreateChatFlowRole(ctx, &vo.ChatFlowRoleCreate{
				WorkflowID: workflowID,
				CreatorID:  ctxutil.MustGetUIDFromCtx(ctx),
			})
			if err != nil {
				return nil, err
			}
		}
	}

	err = GetWorkflowDomainSVC().UpdateMeta(ctx, mustParseInt64(req.GetWorkflowID()), &vo.MetaUpdate{
		Name:         req.Name,
		Desc:         req.Desc,
		IconURI:      req.IconURI,
		WorkflowMode: req.FlowMode,
	})
	if err != nil {
		return nil, err
	}

	safego.Go(ctx, func() {
		err := PublishWorkflowResource(ctx, workflowID, nil, search.Updated, &search.ResourceDocument{
			Name:         req.Name,
			UpdateTimeMS: ptr.Of(time.Now().UnixMilli()),
		})
		if err != nil {
			logs.CtxErrorf(ctx, "publish update workflow resource failed, workflowID: %d, err: %v", workflowID, err)
		}
	})

	return &workflow.UpdateWorkflowMetaResponse{}, nil
}

// CopyWorkflowByWanwu 参考CopyWorkflow，适配了chatflow
func (w *ApplicationService) CopyWorkflowByWanwu(ctx context.Context, req *workflow.CopyWorkflowRequest) (
	resp *workflow.CopyWorkflowResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	spaceID, err := strconv.ParseInt(req.GetSpaceID(), 10, 64)
	if err != nil {
		return nil, err
	}

	if err = checkUserSpace(ctx, ctxutil.MustGetUIDFromCtx(ctx), spaceID); err != nil {
		return nil, err
	}

	workflowID, err := strconv.ParseInt(req.GetWorkflowID(), 10, 64)
	if err != nil {
		return nil, err
	}

	wf, err := w.copyWorkflow(ctx, workflowID, vo.CopyWorkflowPolicy{
		ShouldModifyWorkflowName: true,
	})
	if err != nil {
		return nil, err
	}
	// 对话流的复制需要创建对话
	if wf.Mode == workflow.WorkflowMode_ChatFlow {
		_, err := GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
			AppID:   wf.ID,
			UserID:  wf.CreatorID,
			SpaceID: spaceID,
			Name:    wf.Name,
		})
		if err != nil {
			return nil, err
		}
	}
	return &workflow.CopyWorkflowResponse{
		Data: &workflow.CopyWorkflowData{
			WorkflowID: strconv.FormatInt(wf.ID, 10),
		},
	}, err
}

// ListWorkflowByWanwu 参考ListWorkflow
// 1. size上限 300 -> 99999
// 2. space_id非必须
// 3. login_user_create为true时筛选workflow.CreatorID为当前用户
func (w *ApplicationService) ListWorkflowByWanwu(ctx context.Context, req *workflow.GetWorkFlowListRequest) (
	_ *workflow.GetWorkFlowListResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	option := vo.MetaQuery{}

	if req.GetPage() <= 0 || req.GetSize() <= 0 || req.GetSize() > 99999 {
		return nil, fmt.Errorf("the number of page or size must be greater than 0, and the size must be greater than 0 and less than 99999")
	}
	option.Page = &vo.Page{
		Page: req.GetPage(),
		Size: req.GetSize(),
	}

	userID := ctxutil.MustGetUIDFromCtx(ctx)
	if req.GetSpaceID() != "" {
		spaceID := mustParseInt64(req.GetSpaceID())
		if err := checkUserSpace(ctx, userID, spaceID); err != nil {
			return nil, err
		}
		option.SpaceID = ptr.Of(spaceID)
	}

	if len(req.GetName()) > 0 {
		option.Name = req.Name
	}

	if len(req.GetWorkflowIds()) > 0 {
		ids, err := slices.TransformWithErrorCheck[string, int64](req.GetWorkflowIds(), func(s string) (int64, error) {
			return strconv.ParseInt(s, 10, 64)
		})
		if err != nil {
			return nil, err
		}
		option.IDs = ids
	}

	if req.IsSetFlowMode() && req.GetFlowMode() != workflow.WorkflowMode_All {
		option.Mode = ptr.Of(workflowModel.WorkflowMode(req.GetFlowMode()))
	}

	status := req.GetStatus()
	var qType workflowModel.Locator
	if status == workflow.WorkFlowListStatus_UnPublished {
		option.PublishStatus = ptr.Of(vo.UnPublished)
		qType = workflowModel.FromDraft
	} else if status == workflow.WorkFlowListStatus_HadPublished {
		option.PublishStatus = ptr.Of(vo.HasPublished)
		qType = workflowModel.FromLatestVersion
	}

	wfs, total, err := GetWorkflowDomainSVC().MGet(ctx, &vo.MGetPolicy{
		MetaQuery: option,
		QType:     qType,
		MetaOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	response := &workflow.GetWorkFlowListResponse{
		Data: &workflow.WorkFlowListData{
			AuthList:     make([]*workflow.ResourceAuthInfo, 0),
			WorkflowList: make([]*workflow.Workflow, 0, len(wfs)),
		},
	}

	wf2CreatorID := make(map[int64]string)
	workflowList := make([]*workflow.Workflow, 0, len(wfs))
	for _, w := range wfs {

		if req.GetLoginUserCreate() && w.CreatorID != userID {
			continue
		}

		wf2CreatorID[w.ID] = strconv.FormatInt(w.CreatorID, 10)
		ww := &workflow.Workflow{
			WorkflowID:       strconv.FormatInt(w.ID, 10),
			Name:             w.Name,
			Desc:             w.Desc,
			IconURI:          w.IconURI,
			URL:              w.IconURL,
			CreateTime:       w.CreatedAt.Unix(),
			Type:             w.ContentType,
			SchemaType:       workflow.SchemaType_FDL,
			Tag:              w.Tag,
			TemplateAuthorID: ptr.Of(strconv.FormatInt(w.AuthorID, 10)),
			SpaceID:          ptr.Of(strconv.FormatInt(w.SpaceID, 10)),
			PluginID: func() string {
				if status == workflow.WorkFlowListStatus_UnPublished {
					return "0"
				}
				return strconv.FormatInt(w.ID, 10)
			}(),
			Creator: &workflow.Creator{
				ID:   strconv.FormatInt(w.CreatorID, 10),
				Self: ternary.IFElse[bool](w.CreatorID == ptr.From(ctxutil.GetUIDFromCtx(ctx)), true, false),
			},
			FlowMode: w.Mode,
		}

		if len(req.Checker) > 0 && status == workflow.WorkFlowListStatus_HadPublished {
			ww.CheckResult, err = GetWorkflowDomainSVC().WorkflowSchemaCheck(ctx, w, req.Checker)
			if err != nil {
				return nil, err
			}
		}

		if qType == workflowModel.FromDraft {
			ww.UpdateTime = w.DraftMeta.Timestamp.Unix()
		} else if qType == workflowModel.FromLatestVersion || qType == workflowModel.FromSpecificVersion {
			ww.UpdateTime = w.VersionMeta.VersionCreatedAt.Unix()
		} else if w.UpdatedAt != nil {
			ww.UpdateTime = w.UpdatedAt.Unix()
		}

		startNode := &workflow.Node{
			NodeID:    "100001",
			NodeName:  "start-node",
			NodeParam: &workflow.NodeParam{InputParameters: make([]*workflow.Parameter, 0)},
		}

		for _, in := range w.InputParams {
			param, err := toWorkflowParameter(in)
			if err != nil {
				return nil, err
			}
			startNode.NodeParam.InputParameters = append(startNode.NodeParam.InputParameters, param)
		}

		ww.StartNode = startNode

		auth := &workflow.ResourceAuthInfo{
			WorkflowID: strconv.FormatInt(w.ID, 10),
			UserID:     strconv.FormatInt(w.CreatorID, 10),
			Auth:       &workflow.ResourceActionAuth{CanEdit: true, CanDelete: true, CanCopy: true},
		}
		workflowList = append(workflowList, ww)
		response.Data.AuthList = append(response.Data.AuthList, auth)
	}

	userBasicInfoResponse, err := user.UserApplicationSVC.MGetUserBasicInfo(ctx, &playground.MGetUserBasicInfoRequest{UserIds: slices.Unique(xmaps.Values(wf2CreatorID))})
	if err != nil {
		return nil, err
	}

	for _, w := range workflowList {
		if u, ok := userBasicInfoResponse.UserBasicInfoMap[w.Creator.ID]; ok {
			w.Creator.Name = u.Username
			w.Creator.AvatarURL = u.UserAvatar
		}
	}

	response.Data.WorkflowList = workflowList
	response.Data.Total = total

	return response, nil
}

func (w *ApplicationService) GetWorkFlowSelectByWanwu(ctx context.Context, req *workflow.GetWorkFlowListRequest) (
	_ *workflow.GetWorkFlowListResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	userID := ctxutil.MustGetUIDFromCtx(ctx)
	status := req.GetStatus()
	// 调用ListWorkflowByWanwu需要将status置为nil 查询所有的workflow列表
	// 设置size大小为99999
	req.Status = nil
	sizeValue := int32(99999)
	req.Size = &sizeValue
	resp, err := w.ListWorkflowByWanwu(ctx, req)
	if err != nil {
		return nil, err
	}
	if len(resp.Data.WorkflowList) == 0 {
		return &workflow.GetWorkFlowListResponse{
			Data: &workflow.WorkFlowListData{
				AuthList:     make([]*workflow.ResourceAuthInfo, 0),
				WorkflowList: make([]*workflow.Workflow, 0),
			},
		}, nil
	}
	// 调用BFF的CallBack接口获取发布的workflow IDs
	ret := &cozeWorkflowSelectByWanwuResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"userId": strconv.FormatInt(userID, 10),
			"orgId":  req.GetSpaceID(),
		}).
		SetResult(ret).
		Get(os.Getenv("WANWU_CALLBACK_WORKFLOW_LIST_URL")); err != nil {
		return nil, err
	} else if resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("http request failed with status code: %d", resp.StatusCode())
	}
	// 调用BFF的CallBack接口获取发布的chatflow IDs
	chatRet := &cozeWorkflowSelectByWanwuResp{}
	if resp, err := resty.New().
		R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"userId": strconv.FormatInt(userID, 10),
			"orgId":  req.GetSpaceID(),
		}).
		SetResult(chatRet).
		Get(os.Getenv("WANWU_CALLBACK_CHATFLOW_LIST_URL")); err != nil {
		return nil, err
	} else if resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("http request failed with status code: %d", resp.StatusCode())
	}
	ret.Data.List = append(ret.Data.List, chatRet.Data.List...)
	// 获取已发布的workflow ID集合
	publishedWorkflowIDs := make(map[string]bool)
	for _, wf := range ret.Data.List {
		publishedWorkflowIDs[wf.AppId] = true
	}
	// 过滤WorkflowList
	var filteredWorkflows []*workflow.Workflow
	for _, wf := range resp.Data.WorkflowList {
		isPublished := publishedWorkflowIDs[wf.WorkflowID]
		if status == workflow.WorkFlowListStatus_HadPublished && isPublished {
			filteredWorkflows = append(filteredWorkflows, wf)
		}
		if status == workflow.WorkFlowListStatus_UnPublished && !isPublished {
			wf.PluginID = "0"
			filteredWorkflows = append(filteredWorkflows, wf)
		}
	}
	// 过滤AuthList
	var filteredAuthList []*workflow.ResourceAuthInfo
	for _, wf := range resp.Data.AuthList {
		isPublished := publishedWorkflowIDs[wf.WorkflowID]
		if status == workflow.WorkFlowListStatus_HadPublished && isPublished {
			filteredAuthList = append(filteredAuthList, wf)
		}
		if status == workflow.WorkFlowListStatus_UnPublished && !isPublished {
			filteredAuthList = append(filteredAuthList, wf)
		}
	}
	return &workflow.GetWorkFlowListResponse{
		Data: &workflow.WorkFlowListData{
			AuthList:     filteredAuthList,
			WorkflowList: filteredWorkflows,
			Total:        int64(len(filteredWorkflows)),
		},
	}, nil
}

func (w *ApplicationService) ListWorkFlowOpenAPIV3SchemaByWanwu(ctx context.Context, workflowIDs []string) (
	_ []map[string]any, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowExecuteFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	if len(workflowIDs) == 0 {
		return nil, fmt.Errorf("workflowIDs empty")
	}
	var ids []int64
	for _, workflowID := range workflowIDs {
		ids = append(ids, mustParseInt64(workflowID))
	}

	wfs, _, err := GetWorkflowDomainSVC().MGet(ctx, &vo.MGetPolicy{MetaQuery: vo.MetaQuery{IDs: ids}, QType: workflowModel.FromLatestVersion})
	if err != nil {
		return nil, err
	}
	var rets []map[string]any
	for _, wf := range wfs {
		wfSchema, err := workflowOpenAPIV3Schema(wf)
		if err != nil {
			return nil, err
		}
		rets = append(rets, wfSchema)
	}
	return rets, nil
}

func (w *ApplicationService) GetWorkFlowOpenAPIV3SchemaByWanwu(ctx context.Context, workflowID string) (
	_ map[string]any, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowExecuteFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	wf, err := GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{ID: mustParseInt64(workflowID), QType: workflowModel.FromLatestVersion})
	if err != nil {
		return nil, err
	}
	return workflowOpenAPIV3Schema(wf)
}

func workflowOpenAPIV3Schema(wf *entity.Workflow) (map[string]any, error) {
	inputSchema, err := workflowParamsToSchema(wf.InputParams)
	if err != nil {
		return nil, fmt.Errorf("workflow(%v) input params to openapi v3 schema err: %v", wf.ID, err)
	}
	outputSchema, err := workflowParamsToSchema(wf.OutputParams)
	if err != nil {
		return nil, fmt.Errorf("workflow(%v) output params to openapi v3 schema err: %v", wf.ID, err)
	}

	var input any = inputSchema
	if inputSchema == nil {
		input = map[string]string{
			"type": "object",
		}
	}
	var output any = outputSchema
	if outputSchema == nil {
		output = map[string]string{
			"type": "object",
		}
	}

	return map[string]any{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":       wf.Name,
			"version":     "1.0.0",
			"description": wf.Desc,
		},
		"servers": []map[string]string{
			{
				"url": "http://workflow-wanwu:8999/v1",
			},
		},
		"paths": map[string]any{
			path.Join("/workflow", strconv.Itoa(int(wf.ID)), "/run_by_wanwu"): map[string]any{
				"post": map[string]any{
					"summary":     wf.Name,
					"operationId": "action_" + wf.Name,
					"description": wf.Desc,
					"parameters": []map[string]any{
						{
							"in":   "header",
							"name": "Content-Type",
							"schema": map[string]string{
								"type":    "string",
								"example": "application/json",
							},
							"required": true,
						},
					},
					"requestBody": map[string]any{
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": input,
							},
						},
					},
					"responses": map[string]any{
						"200": map[string]any{
							"description": "请求成功时的结果",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": output,
								},
							},
						},
						"default": map[string]any{
							"description": "请求失败时的错误信息",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]string{
										"type": "object",
									},
								},
							},
						},
					},
				},
			},
		},
	}, nil
}

func workflowParamsToSchema(params []*vo.NamedTypeInfo) (*openapi3.Schema, error) {
	var paramsMap map[string]*schema.ParameterInfo
	for _, param := range params {
		paramInfo, err := param.ToParameterInfo()
		if err != nil {
			return nil, err
		}
		if paramsMap == nil {
			paramsMap = make(map[string]*schema.ParameterInfo)
		}
		paramsMap[param.Name] = paramInfo
	}
	return schema.NewParamsOneOfByParams(paramsMap).ToOpenAPIV3()
}

// OpenAPIRunByWanwu 参考OpenAPIRun
// 0. FIXME 智能体运行该接口，不会在header中带userId、orgId，需要在该方法中设置ctxcache
// 1. 去掉api auth、user check等业务逻辑
// 2. 去掉appID、agentID、connectorID等业务逻辑
// 3. 将必须publish才能执行的workflow，改为可以执行draft
func (w *ApplicationService) OpenAPIRunByWanwu(ctx context.Context, workflowID string, req *workflow.OpenAPIRunFlowRequest) (
	_ *workflow.OpenAPIRunFlowResponse, _ vo.TerminatePlan, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowExecuteFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	// apiKeyInfo := ctxutil.GetApiAuthFromCtx(ctx)
	// userID := apiKeyInfo.UserID

	parameters := make(map[string]any)
	if req.Parameters != nil {
		err := sonic.UnmarshalString(*req.Parameters, &parameters)
		if err != nil {
			return nil, "", vo.WrapError(errno.ErrInvalidParameter, err)
		}
	}

	meta, err := GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{
		ID:       mustParseInt64(workflowID),
		MetaOnly: true,
		QType:    workflowModel.FromLatestVersion,
	})
	if err != nil {
		return nil, "", err
	}

	// 设置ctxcache
	if ctxutil.GetUserSessionFromCtx(ctx) == nil {
		ctxcache.Store(ctx, consts.SessionDataKeyInCtx, &user_entity.Session{
			UserID: meta.CreatorID,
			Locale: string(i18n.GetLocale(ctx)),
		})
		ctxcache.Store(ctx, "X-Org-Id", strconv.Itoa(int(meta.SpaceID)))
	}

	if meta.LatestPublishedVersion == nil {
		return nil, "", vo.NewError(errno.ErrWorkflowNotPublished)
	}

	// if err = checkUserSpace(ctx, userID, meta.SpaceID); err != nil {
	// 	return nil, err
	// }

	// var appID, agentID *int64
	// if req.IsSetAppID() {
	// 	appID = ptr.Of(mustParseInt64(req.GetAppID()))
	// } else if req.IsSetProjectID() {
	// 	appID = ptr.Of(mustParseInt64(req.GetProjectID()))
	// }
	// if req.IsSetBotID() {
	// 	agentID = ptr.Of(mustParseInt64(req.GetBotID()))
	// }

	// var connectorID int64
	// if req.IsSetConnectorID() {
	// 	connectorID = mustParseInt64(req.GetConnectorID())
	// }

	// if connectorID != consts.WebSDKConnectorID {
	// 	connectorID = apiKeyInfo.ConnectorID
	// }

	exeCfg := workflowModel.ExecuteConfig{
		ID:       meta.ID,
		From:     workflowModel.FromLatestVersion,
		Version:  *meta.LatestPublishedVersion,
		Operator: meta.CreatorID,
		Mode:     workflowModel.ExecuteModeRelease,
		// AppID:         appID,
		// AgentID:       agentID,
		ConnectorID:   consts.WebSDKConnectorID,
		ConnectorUID:  strconv.FormatInt(meta.CreatorID, 10),
		InputFailFast: true,
		BizType:       workflowModel.BizTypeWorkflow,
	}

	// if exeCfg.AppID != nil && exeCfg.AgentID != nil {
	// 	return nil, errors.New("project_id and bot_id cannot be set at the same time")
	// }

	if req.GetIsAsync() {
		exeCfg.SyncPattern = workflowModel.SyncPatternAsync
		exeCfg.TaskType = workflowModel.TaskTypeBackground
		exeID, err := GetWorkflowDomainSVC().AsyncExecute(ctx, exeCfg, parameters)
		if err != nil {
			return nil, "", err
		}

		return &workflow.OpenAPIRunFlowResponse{
			ExecuteID: ptr.Of(strconv.FormatInt(exeID, 10)),
			DebugUrl:  ptr.Of(debugutil.GetWorkflowDebugURL(ctx, meta.ID, meta.SpaceID, exeID)),
		}, "", nil
	}

	exeCfg.SyncPattern = workflowModel.SyncPatternSync
	exeCfg.TaskType = workflowModel.TaskTypeForeground
	wfExe, tPlan, err := GetWorkflowDomainSVC().SyncExecute(ctx, exeCfg, parameters)
	if err != nil {
		return nil, "", err
	}

	if wfExe.Status == entity.WorkflowInterrupted {
		return nil, "", vo.NewError(errno.ErrInterruptNotSupported)
	}

	data := wfExe.Output
	return &workflow.OpenAPIRunFlowResponse{
		Data:      data,
		ExecuteID: ptr.Of(strconv.FormatInt(wfExe.ID, 10)),
		DebugUrl:  ptr.Of(debugutil.GetWorkflowDebugURL(ctx, meta.ID, wfExe.SpaceID, wfExe.ID)),
		Token:     ptr.Of(wfExe.TokenInfo.InputTokens + wfExe.TokenInfo.OutputTokens),
		Cost:      ptr.Of("0.00000"),
	}, tPlan, nil
}

// GetPlaygroundPluginListByWanwu 参考GetPlaygroundPluginList
func (w *ApplicationService) GetPlaygroundPluginListByWanwu(ctx context.Context, req *pluginAPI.GetPlaygroundPluginListRequest) (
	resp *pluginAPI.GetPlaygroundPluginListResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	currentUser := ctxutil.MustGetUIDFromCtx(ctx)
	if err = checkUserSpace(ctx, currentUser, req.GetSpaceID()); err != nil {
		return nil, err
	}

	var (
		toolIDs []int64
		wfs     []*entity.Workflow
	)
	if len(req.GetPluginIds()) > 0 {
		toolIDs, err = slices.TransformWithErrorCheck(req.GetPluginIds(), func(a string) (int64, error) {
			return strconv.ParseInt(a, 10, 64)
		})
		if err != nil {
			return nil, err
		}
		// 修改QType从Draft中查找
		wfs, _, err = GetWorkflowDomainSVC().MGet(ctx, &vo.MGetPolicy{
			MetaQuery: vo.MetaQuery{
				IDs:     toolIDs,
				SpaceID: ptr.Of(req.GetSpaceID()),
			},
			QType: workflowModel.FromDraft,
		})
	} else if req.GetPage() > 0 && req.GetSize() > 0 {
		wfs, _, err = GetWorkflowDomainSVC().MGet(ctx, &vo.MGetPolicy{
			MetaQuery: vo.MetaQuery{
				Page: &vo.Page{
					Size: req.GetSize(),
					Page: req.GetPage(),
				},
				SpaceID:       ptr.Of(req.GetSpaceID()),
				PublishStatus: ptr.Of(vo.HasPublished),
			},
			QType: workflowModel.FromLatestVersion,
		})
	}

	if err != nil {
		return nil, err
	}

	pluginInfoList := make([]*common.PluginInfoForPlayground, 0)
	for _, wf := range wfs {
		pInfo := &common.PluginInfoForPlayground{
			ID:           strconv.FormatInt(wf.ID, 10),
			Name:         wf.Name,
			PluginIcon:   wf.IconURL,
			DescForHuman: wf.Desc,
			Creator: &common.Creator{
				Self: wf.CreatorID == currentUser,
			},
			PluginType: common.PluginType_WORKFLOW,
			//VersionName: wf.VersionMeta.Version,
			//CreateTime:  strconv.FormatInt(wf.CreatedAt.Unix(), 10),
			//UpdateTime:  strconv.FormatInt(wf.VersionCreatedAt.Unix(), 10),
		}

		pluginApi := &common.PluginApi{
			APIID:    strconv.FormatInt(wf.ID, 10),
			Name:     wf.Name,
			Desc:     wf.Desc,
			PluginID: strconv.FormatInt(wf.ID, 10),
		}
		pluginApi.Parameters, err = slices.TransformWithErrorCheck(wf.InputParams, toPluginParameter)
		if err != nil {
			return nil, err
		}

		pInfo.PluginApis = []*common.PluginApi{pluginApi}
		pluginInfoList = append(pluginInfoList, pInfo)
	}

	return &pluginAPI.GetPlaygroundPluginListResponse{
		Data: &common.GetPlaygroundPluginListData{
			PluginList: pluginInfoList,
			Total:      int32(len(pluginInfoList)),
		},
	}, nil
}

// LatestVersionRunByWanwu 参考 TestRun
// 执行已发布的工作流
func (w *ApplicationService) LatestVersionRunByWanwu(ctx context.Context, req *workflow.WorkFlowTestRunRequest) (_ *workflow.WorkFlowTestRunResponse, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowExecuteFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	uID := ctxutil.MustGetUIDFromCtx(ctx)

	if err := checkUserSpace(ctx, uID, mustParseInt64(req.GetSpaceID())); err != nil {
		return nil, err
	}

	var appID, agentID *int64
	if req.IsSetProjectID() {
		appID = ptr.Of(mustParseInt64(req.GetProjectID()))
	}
	if req.IsSetBotID() {
		agentID = ptr.Of(mustParseInt64(req.GetBotID()))
	}

	exeCfg := workflowModel.ExecuteConfig{
		ID:           mustParseInt64(req.GetWorkflowID()),
		From:         workflowModel.FromLatestVersion,
		CommitID:     req.GetCommitID(),
		Operator:     uID,
		Mode:         workflowModel.ExecuteModeDebug,
		AppID:        appID,
		AgentID:      agentID,
		ConnectorID:  consts.CozeConnectorID,
		ConnectorUID: strconv.FormatInt(uID, 10),
		TaskType:     workflowModel.TaskTypeForeground,
		SyncPattern:  workflowModel.SyncPatternAsync,
		BizType:      workflowModel.BizTypeWorkflow,
		Cancellable:  true,
	}

	if exeCfg.AppID != nil && exeCfg.AgentID != nil {
		return nil, errors.New("project_id and bot_id cannot be set at the same time")
	}

	exeID, err := GetWorkflowDomainSVC().AsyncExecute(ctx, exeCfg, maps.ToAnyValue(req.Input))
	if err != nil {
		return nil, err
	}

	return &workflow.WorkFlowTestRunResponse{
		Data: &workflow.WorkFlowTestRunData{
			WorkflowID: req.WorkflowID,
			ExecuteID:  fmt.Sprintf("%d", exeID),
		},
	}, nil
}

// GetWorkflowVersionListByWanwu 获取工作流的所有版本列表
func (w *ApplicationService) GetWorkflowVersionListByWanwu(ctx context.Context, workflowID string) (
	_ *GetWorkflowVersionListResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}
	}()
	versionInfos, err := GetWorkflowDomainSVC().GetWorkflowVersionListByWanwu(ctx, mustParseInt64(workflowID))
	if err != nil {
		return nil, err
	}
	response := GetWorkflowVersionListResponse{
		Data: &WorkflowVersionListData{
			WorkflowID:  workflowID,
			VersionList: make([]*WorkflowVersion, 0, len(versionInfos)),
			Total:       int32(len(versionInfos)),
		},
	}
	for _, versionInfo := range versionInfos {
		wv := &WorkflowVersion{
			Version:            versionInfo.Version,
			VersionDescription: versionInfo.VersionDescription,
			CreatedAt:          versionInfo.VersionCreatedAt.UnixMilli(),
			CommitID:           workflowID + "_" + versionInfo.Version,
			Type:               workflow.OperateType_PublishOperate,
		}
		response.Data.VersionList = append(response.Data.VersionList, wv)
	}

	return &response, nil
}

func (w *ApplicationService) UpdateWorkflowVersionDescriptionByWanwu(ctx context.Context, req *UpdateWorkflowVersionDescriptionRequest) (_ *workflow.UpdateWorkflowMetaResponse, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}
	}()
	err = GetWorkflowDomainSVC().UpdateWorkflowVersionDescriptionByWanwu(ctx, mustParseInt64(req.WorkflowID), req.VersionDescription)
	if err != nil {
		return nil, err
	}
	return &workflow.UpdateWorkflowMetaResponse{}, nil
}

func (w *ApplicationService) RollbackWorkflowVersionByWanwu(ctx context.Context, req *RollbackWorkflowVersionRequest) (
	_ *workflow.SaveWorkflowResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	var version string
	if req.CommitID != "" {
		if idx := strings.LastIndex(req.CommitID, "_"); idx != -1 && idx < len(req.CommitID)-1 {
			version = req.CommitID[idx+1:]
		} else {
			return nil, fmt.Errorf("invalid commit_id format: %s", req.CommitID)
		}
	}
	policy := &vo.GetPolicy{
		ID:      mustParseInt64(req.WorkflowID),
		QType:   ternary.IFElse(len(version) > 0, workflowModel.FromSpecificVersion, workflowModel.FromDraft),
		Version: version,
	}

	wfEntity, err := GetWorkflowDomainSVC().Get(ctx, policy)
	if err != nil {
		return nil, err
	}
	if err := GetWorkflowDomainSVC().Save(ctx, mustParseInt64(req.WorkflowID), wfEntity.Canvas); err != nil {
		return nil, err
	}

	return &workflow.SaveWorkflowResponse{
		Data: &workflow.SaveWorkflowData{},
	}, nil
}

func (w *ApplicationService) ExportWorkFlowByWanwu(ctx context.Context, req *ExportWorkflowRequest) (
	_ *workflow.GetHistorySchemaResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	var qType workflowModel.Locator
	if req.Version != "" {
		qType = workflowModel.FromSpecificVersion
	} else {
		qType = workflowModel.FromDraft
	}

	policy := &vo.GetPolicy{
		ID:      mustParseInt64(req.WorkflowID),
		QType:   qType,
		Version: req.Version,
	}

	wfEntity, err := GetWorkflowDomainSVC().Get(ctx, policy)
	if err != nil {
		return nil, err
	}
	return &workflow.GetHistorySchemaResponse{
		Data: &workflow.GetHistorySchemaData{
			Name:     wfEntity.Name,
			Describe: wfEntity.Desc,
			Schema:   wfEntity.Canvas,
		},
	}, nil
}

func (w *ApplicationService) GetWorkflowVersionSchema(ctx context.Context, req *workflow.GetHistorySchemaRequest) (
	resp *workflow.GetHistorySchemaResponse, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	if req.GetSpaceID() == "" {
		return nil, vo.WrapError(errno.ErrInvalidParameter, fmt.Errorf("space_id is required"))
	}
	if req.GetWorkflowID() == "" {
		return nil, vo.WrapError(errno.ErrInvalidParameter, fmt.Errorf("workflow_id is required"))
	}

	var version string
	if req.CommitID != nil {
		commitID := *req.CommitID
		if idx := strings.LastIndex(commitID, "_"); idx != -1 && idx < len(commitID)-1 {
			version = commitID[idx+1:]
		} else {
			return nil, fmt.Errorf("invalid commit_id format: %s", req.CommitID)
		}
	}
	policy := &vo.GetPolicy{
		ID:      mustParseInt64(req.GetWorkflowID()),
		QType:   workflowModel.FromSpecificVersion,
		Version: version,
	}

	wfEntity, err := GetWorkflowDomainSVC().Get(ctx, policy)
	return &workflow.GetHistorySchemaResponse{
		Data: &workflow.GetHistorySchemaData{
			Name:       wfEntity.Name,
			Describe:   wfEntity.Desc,
			URL:        wfEntity.IconURL,
			Schema:     wfEntity.Canvas,
			FlowMode:   wfEntity.Mode,
			WorkflowID: req.GetWorkflowID(),
			CommitID:   req.GetWorkflowID() + "_" + version,
		},
	}, nil
}

// GetCanvasInfoByWanwu 参考GetCanvasInfo
func (w *ApplicationService) GetCanvasInfoByWanwu(ctx context.Context, req *ExportWorkflowRequest) (
	_ *workflow.GetCanvasInfoResponse, err error,
) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}
		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	switch req.Type {
	case workflowModel.FromDraft:
		return w.GetCanvasInfo(ctx, &workflow.GetCanvasInfoRequest{
			SpaceID:    req.SpaceID,
			WorkflowID: ptr.Of(req.WorkflowID),
		})

	case workflowModel.FromSpecificVersion:
		return nil, fmt.Errorf("invalid type")
	}

	wf, err := GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{
		ID:    mustParseInt64(req.WorkflowID),
		QType: workflowModel.FromLatestVersion,
	})
	if err != nil {
		return nil, err
	}
	var devStatus workflow.WorkFlowDevStatus
	var vcsType workflow.VCSCanvasType
	pluginID := "0"
	updateTime := time.Time{}

	canvasData := &workflow.CanvasData{
		Workflow: &workflow.Workflow{
			WorkflowID:       strconv.FormatInt(wf.ID, 10),
			Name:             wf.Name,
			Desc:             wf.Desc,
			URL:              wf.IconURL,
			IconURI:          wf.IconURI,
			Status:           devStatus,
			Type:             wf.ContentType,
			CreateTime:       wf.CreatedAt.Unix(),
			UpdateTime:       updateTime.Unix(),
			Tag:              wf.Tag,
			TemplateAuthorID: ternary.IFElse(wf.AuthorID > 0, ptr.Of(strconv.FormatInt(wf.AuthorID, 10)), nil),
			SpaceID:          ptr.Of(strconv.FormatInt(wf.SpaceID, 10)),
			SchemaJSON:       ptr.Of(wf.Canvas),
			Creator: &workflow.Creator{
				ID:   strconv.FormatInt(wf.CreatorID, 10),
				Self: ternary.IFElse[bool](wf.CreatorID == ptr.From(ctxutil.GetUIDFromCtx(ctx)), true, false),
			},
			FlowMode:         wf.Mode,
			ProjectID:        i64PtrToStringPtr(wf.AppID),
			PersistenceModel: workflow.PersistenceModel_VCS, // the front-end validation logic, this field returns VCS, developers don't need to pay attention
			PluginID:         pluginID,
		},
		VcsData: &workflow.VCSCanvasData{
			CanEdit:        true,
			SubmitCommitID: wf.CommitID,
			DraftCommitID:  wf.CommitID,
			Type:           vcsType,
		},
		WorkflowVersion: wf.LatestPublishedVersion,
	}

	return &workflow.GetCanvasInfoResponse{
		Data: canvasData,
	}, nil
}

// -- internal ---
type cozeWorkflowSelectByWanwuResp struct {
	Code int64 `json:"code"`
	Data struct {
		List []appId `json:"list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type appId struct {
	AppId string `json:"appId"` // 应用id
}

type GetWorkflowVersionListResponse struct {
	Data     *WorkflowVersionListData `thrift:"data,1,required" form:"data,required" json:"data,required" query:"data,required"`
	Code     int64                    `thrift:"code,253,required" form:"code,required" json:"code,required" query:"code,required"`
	Msg      string                   `thrift:"msg,254,required" form:"msg,required" json:"msg,required" query:"msg,required"`
	BaseResp *base.BaseResp           `thrift:"BaseResp,255,required" form:"BaseResp,required" json:"BaseResp,required" query:"BaseResp,required"`
}

type WorkflowVersionListData struct {
	WorkflowID  string             `thrift:"workflow_id,1" form:"workflow_id" json:"workflow_id" query:"workflow_id"`
	VersionList []*WorkflowVersion `thrift:"version_list,2" form:"version_list" json:"version_list" query:"version_list"`
	Total       int32              `thrift:"total,3" form:"total" json:"total" query:"total"`
}

type WorkflowVersion struct {
	Version            string               `thrift:"version,1" form:"version" json:"version" query:"version"`
	VersionDescription string               `thrift:"version_description,2" form:"version_description" json:"version_description" query:"version_description"`
	CreatedAt          int64                `thrift:"created_at,3" form:"created_at" json:"created_at" query:"created_at"`
	CommitID           string               `thrift:"commit_id,3,optional" form:"commit_id" json:"commit_id,omitempty" query:"commit_id"`
	Type               workflow.OperateType `thrift:"type,4,required" form:"type,required" json:"type,required" query:"type,required"`
}

type UpdateWorkflowVersionDescriptionRequest struct {
	WorkflowID         string `thrift:"workflow_id,1,required" form:"workflow_id,required" json:"workflow_id,required" query:"workflow_id,required"`
	VersionDescription string `thrift:"version_description,8,optional" form:"version_description" json:"version_description,omitempty" query:"version_description"`
}

type RollbackWorkflowVersionRequest struct {
	WorkflowID string  `thrift:"workflow_id,1,required" form:"workflow_id,required" json:"workflow_id,required" query:"workflow_id,required"`
	Version    *string `thrift:"version,6,optional" form:"version" json:"version,omitempty" query:"version"`
	CommitID   string  `thrift:"commit_id,3,optional" form:"commit_id" json:"commit_id,omitempty" query:"commit_id"`
}

type GetWorkflowRequest struct {
	WorkflowID string `thrift:"workflow_id,1,required" form:"workflow_id,required" json:"workflow_id,required" query:"workflow_id,required"`
}

type RunWorkFlowLatestVersionByWanwuReq struct {
	WorkflowID string         `json:"workflow_id,required" validate:"required"`
	Input      map[string]any `json:"input,omitempty"`
}

type GetWorkflowLatestVersionCanvasResponse struct {
	Data     *WorkflowLatestVersionCanvasInfo `thrift:"data,1,required" json:"data,required"`
	Code     int64                            `thrift:"code,253,required" json:"code,required"`
	Msg      string                           `thrift:"msg,254,required" json:"msg,required"`
	BaseResp *base.BaseResp                   `thrift:"BaseResp,255,required" json:"BaseResp,required"`
}

type WorkflowLatestVersionCanvasInfo struct {
	Name               string `thrift:"name,1" json:"name"`
	Describe           string `thrift:"describe,2" json:"describe"`
	WorkflowID         string `thrift:"workflow_id,3" json:"workflow_id"`
	Schema             string `thrift:"schema,4" json:"schema"`
	Version            string `thrift:"version,5" json:"version"`
	VersionDescription string `thrift:"version_description,6" json:"version_description"`
}

type ExportWorkflowRequest struct {
	WorkflowID string `thrift:"workflow_id,1,required" form:"workflow_id,required" json:"workflow_id,required" query:"workflow_id,required"`
	SpaceID    string `form:"space_id,required" json:"space_id" query:"space_id,required"`
	Version    string `thrift:"version,6,optional" form:"version" json:"version,omitempty" query:"version"`
	Type       model.Locator
}
