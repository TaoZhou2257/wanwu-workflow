package workflow

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"

	workflowModel "github.com/coze-dev/coze-studio/backend/api/model/crossdomain/workflow"
	"github.com/coze-dev/coze-studio/backend/api/model/playground"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	"github.com/coze-dev/coze-studio/backend/application/user"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/slices"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ternary"
	"github.com/coze-dev/coze-studio/backend/pkg/safego"
	"github.com/coze-dev/coze-studio/backend/types/errno"
	xmaps "golang.org/x/exp/maps"
)

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
