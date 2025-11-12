package workflow

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	crossconversation "github.com/coze-dev/coze-studio/backend/crossdomain/conversation"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ternary"
	"github.com/coze-dev/coze-studio/backend/pkg/safego"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

func (w *ApplicationService) OpenAPICreateConversationByWanwu(ctx context.Context, req *workflow.CreateConversationRequest) (resp *workflow.CreateConversationResponse, err error) {

	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}
		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	var (
		appID      = mustParseInt64(req.GetAppID())
		apiKeyInfo = ctxutil.GetApiAuthFromCtx(ctx)
		userID     = apiKeyInfo.UserID
		env        = ternary.IFElse(req.GetDraftMode(), vo.Draft, vo.Online)
		cID        int64
		//spaceID = mustParseInt64(req.GetSpaceID())
		//_       = spaceID
	)

	// todo  check permission

	if !req.GetGetOrCreate() {
		cID, err = GetWorkflowDomainSVC().UpdateConversation(ctx, env, appID, req.GetConnectorId(), userID, req.GetConversationMame())
	} else {
		var tplExisted, dcExisted bool
		var tplErr, dcErr error
		var wg sync.WaitGroup
		wg.Add(2)

		safego.Go(ctx, func() {
			defer wg.Done()
			_, tplExisted, tplErr = GetWorkflowDomainSVC().GetTemplateByName(ctx, env, appID, req.GetConversationMame())
		})

		safego.Go(ctx, func() {
			defer wg.Done()
			_, dcExisted, dcErr = GetWorkflowDomainSVC().GetDynamicConversationByName(ctx, env, appID, req.GetConnectorId(), userID, req.GetConversationMame())
		})

		wg.Wait()

		if tplErr != nil {
			return nil, tplErr
		}
		if dcErr != nil {
			return nil, dcErr
		}
		if !tplExisted && !dcExisted {
			// 尝试创建
			createdTemplate, err := GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
				AppID:   appID,
				Name:    "Default",
				UserID:  userID,
				SpaceID: 1,
			})
			if err != nil {
				// 创建失败
				return nil, err
			}
			if createdTemplate == 0 {
				// 创建返回空，说明创建未成功
				return &workflow.CreateConversationResponse{
					Code: errno.ErrConversationNotFoundForOperation,
					Msg:  "Failed to create conversation. Please try again.",
				}, nil
			}
		}
		cID, _, err = GetWorkflowDomainSVC().GetOrCreateConversation(ctx, env, appID, req.GetConnectorId(), userID, req.GetConversationMame())

	}
	if err != nil {
		return nil, err
	}

	cInfo, err := crossconversation.DefaultSVC().GetByID(ctx, cID)
	if err != nil {
		return nil, err
	}
	return &workflow.CreateConversationResponse{
		ConversationData: &workflow.ConversationData{
			Id:            cID,
			LastSectionID: ptr.Of(cInfo.SectionID),
		},
	}, nil
}
