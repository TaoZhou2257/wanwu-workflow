package workflow

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	crossconversation "github.com/coze-dev/coze-studio/backend/crossdomain/conversation"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ternary"
	"github.com/coze-dev/coze-studio/backend/pkg/safego"
	"github.com/coze-dev/coze-studio/backend/types/errno"
)

// CreateApplicationConversationDefByWanwu 参考CreateApplicationConversationDef
func (w *ApplicationService) CreateApplicationConversationDefByWanwu(ctx context.Context, req *workflow.CreateProjectConversationDefRequest) (resp *workflow.CreateProjectConversationDefResponse, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrConversationOfAppOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()
	// 取消获取appID，返回随机生成的newAppId（放在uniqueID字段中）作为app-service的applicationId
	var (
		spaceID = mustParseInt64(req.GetSpaceID())
		//appID   = mustParseInt64(req.GetProjectID())
		userID = ctxutil.MustGetUIDFromCtx(ctx)
	)

	if err := checkUserSpace(ctx, userID, spaceID); err != nil {
		return nil, err
	}
	newAppID, _ := w.IDGenerator.GenID(ctx)
	// 应用广场默认创建conversation template，表现为用户首次进入应用广场对话流，会有一个默认的conversation
	_, err = GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
		AppID:   newAppID,
		SpaceID: spaceID,
		Name:    req.GetConversationName(),
		UserID:  userID,
	})
	if err != nil {
		return nil, err
	}

	return &workflow.CreateProjectConversationDefResponse{
		UniqueID: strconv.FormatInt(newAppID, 10),
		SpaceID:  req.GetSpaceID(),
	}, err
}

// OpenAPICreateConversationByWanwu 参考OpenAPICreateConversation
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
		// 先从req中获取spaceID（如果是前端调用，req中不会传orgId要从header中获取）
		spaceID = req.GetSpaceID()

		templateId int64
		t          *entity.ConversationTemplate
	)

	// todo  check permission

	if spaceID == "" {
		spaceID, _ = ctxcache.Get[string](ctx, "X-Org-Id")
	}

	if !req.GetGetOrCreate() {
		cID, err = GetWorkflowDomainSVC().UpdateConversation(ctx, env, appID, req.GetConnectorId(), userID, req.GetConversationMame())
	} else {
		var tplExisted, dcExisted bool
		var tplErr, dcErr error
		var wg sync.WaitGroup
		wg.Add(2)

		safego.Go(ctx, func() {
			defer wg.Done()
			t, tplExisted, tplErr = GetWorkflowDomainSVC().GetTemplateByName(ctx, env, appID, req.GetConversationMame())
			if tplExisted {
				// 需要给前端返回templateId
				templateId = t.TemplateID
			}
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
			// 应用广场用户新建会话，需要先创建conversation template
			templateId, err = GetWorkflowDomainSVC().CreateDraftConversationTemplate(ctx, &vo.CreateConversationTemplateMeta{
				AppID:   appID,
				UserID:  userID,
				SpaceID: mustParseInt64(spaceID),
				Name:    *req.ConversationMame,
			})
			if err != nil {
				return &workflow.CreateConversationResponse{
					Code: errno.ErrConversationNotFoundForOperation,
					Msg:  fmt.Sprintf("Conversation not found. Please create a conversation before attempting to perform any related operations. (%v)", err),
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
			MetaData: map[string]string{
				// 将templateId返回给前端
				"uniqueId": strconv.Itoa(int(templateId)),
				// 将appId返回给bff
				"appId": req.GetAppID(),
			},
		},
	}, nil
}
