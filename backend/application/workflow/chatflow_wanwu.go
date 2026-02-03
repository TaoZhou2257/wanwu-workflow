package workflow

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"sync"

	"github.com/cloudwego/eino/schema"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	crossagentrun "github.com/coze-dev/coze-studio/backend/crossdomain/agentrun"
	crossconversation "github.com/coze-dev/coze-studio/backend/crossdomain/conversation"
	crossmessage "github.com/coze-dev/coze-studio/backend/crossdomain/message"
	message "github.com/coze-dev/coze-studio/backend/crossdomain/message/model"
	workflowModel "github.com/coze-dev/coze-studio/backend/crossdomain/workflow/model"
	agententity "github.com/coze-dev/coze-studio/backend/domain/conversation/agentrun/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
	"github.com/coze-dev/coze-studio/backend/pkg/errorx"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ptr"
	"github.com/coze-dev/coze-studio/backend/pkg/lang/ternary"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/safego"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
	"github.com/coze-dev/coze-studio/backend/pkg/taskgroup"
	"github.com/coze-dev/coze-studio/backend/types/consts"
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
			// 只从drafttemplate中获取
			t, tplExisted, tplErr = GetWorkflowDomainSVC().GetTemplateByName(ctx, vo.Draft, appID, req.GetConversationMame())
			if tplExisted {
				// 需要给前端返回templateId
				templateId = t.TemplateID
			}
		})

		safego.Go(ctx, func() {
			defer wg.Done()
			_, dcExisted, dcErr = GetWorkflowDomainSVC().GetDynamicConversationByName(ctx, vo.Draft, appID, req.GetConnectorId(), userID, req.GetConversationMame())
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

// OpenAPIChatFlowRunByWanwu 参考OpenAPIChatFlowRun
func (w *ApplicationService) OpenAPIChatFlowRunByWanwu(ctx context.Context, req *workflow.ChatFlowRunRequest) (
	_ *schema.StreamReader[[]*workflow.ChatFlowRunResponse], err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			err = safego.NewPanicErr(panicErr, debug.Stack())
		}

		if err != nil {
			err = vo.WrapIfNeeded(errno.ErrWorkflowOperationFail, err, errorx.KV("cause", vo.UnwrapRootErr(err).Error()))
		}
	}()

	if len(req.GetAdditionalMessages()) == 0 {
		return nil, fmt.Errorf("additional_messages is requird")
	}

	messages := req.GetAdditionalMessages()

	lastUserMessage := messages[len(req.GetAdditionalMessages())-1]
	if lastUserMessage.Role != userRole {
		return nil, errors.New("the role of the last day message must be user")
	}

	var parameters = make(map[string]any)
	if len(req.GetParameters()) > 0 {
		err := sonic.UnmarshalString(req.GetParameters(), &parameters)
		if err != nil {
			return nil, err
		}
		// 防止 JSON "null" 导致 map 变成 nil
		if parameters == nil {
			parameters = make(map[string]any)
		}
	}

	var (
		workflowID     = mustParseInt64(req.GetWorkflowID())
		isDebug        = req.GetExecuteMode() == "DEBUG"
		appID, agentID *int64
		bizID          int64
		conversationID int64
		sectionID      int64
		version        string
		locator        workflowModel.Locator
		apiKeyInfo     = ctxutil.GetApiAuthFromCtx(ctx)
		userID         = apiKeyInfo.UserID
		connectorID    int64
	)
	if len(req.GetConnectorID()) == 0 {
		connectorID = ternary.IFElse(isDebug, consts.CozeConnectorID, apiKeyInfo.ConnectorID)
	} else {
		connectorID = mustParseInt64(req.GetConnectorID())
	}

	if req.IsSetAppID() {
		appID = ptr.Of(mustParseInt64(req.GetAppID()))
		bizID = mustParseInt64(req.GetAppID())
	}
	if req.IsSetBotID() {
		agentID = ptr.Of(mustParseInt64(req.GetBotID()))
		bizID = mustParseInt64(req.GetBotID())
	}

	if appID != nil && agentID != nil {
		return nil, errors.New("project_id and bot_id cannot be set at the same time")
	}

	if isDebug {
		locator = workflowModel.FromDraft
	} else {
		meta, err := GetWorkflowDomainSVC().Get(ctx, &vo.GetPolicy{
			ID:       workflowID,
			MetaOnly: true,
		})
		if err != nil {
			return nil, err
		}

		if meta.LatestPublishedVersion == nil {
			return nil, vo.NewError(errno.ErrWorkflowNotPublished)
		}
		if req.IsSetVersion() {
			version = req.GetVersion()
			locator = workflowModel.FromSpecificVersion
		} else {
			version = meta.GetLatestVersion()
			locator = workflowModel.FromLatestVersion
		}
	}

	if req.IsSetConversationID() && !req.IsSetBotID() {
		conversationID = mustParseInt64(req.GetConversationID())
		cInfo, err := crossconversation.DefaultSVC().GetByID(ctx, conversationID)
		if err != nil {
			return nil, err
		}
		sectionID = cInfo.SectionID

		// only trust the conversation name under the app
		conversationName, existed, err := GetWorkflowDomainSVC().GetConversationNameByID(ctx, ternary.IFElse(isDebug, vo.Draft, vo.Online), bizID, connectorID, conversationID)
		if err != nil {
			return nil, err
		}
		if !existed {
			return nil, fmt.Errorf("conversation not found")
		}
		parameters[vo.ConversationNameKey] = conversationName
	} else if req.IsSetConversationID() && req.IsSetBotID() {
		parameters[vo.ConversationNameKey] = "Default"
		conversationID = mustParseInt64(req.GetConversationID())
		cInfo, err := crossconversation.DefaultSVC().GetByID(ctx, conversationID)
		if err != nil {
			return nil, err
		}
		sectionID = cInfo.SectionID
	} else {
		conversationName, ok := parameters[vo.ConversationNameKey].(string)
		if !ok {
			return nil, fmt.Errorf("conversation name is requried")
		}
		cID, sID, err := GetWorkflowDomainSVC().GetOrCreateConversation(ctx, ternary.IFElse(isDebug, vo.Draft, vo.Online), bizID, connectorID, userID, conversationName)
		if err != nil {
			return nil, err
		}
		conversationID = cID
		sectionID = sID
	}

	runRecord, err := crossagentrun.DefaultSVC().Create(ctx, &agententity.AgentRunMeta{
		AgentID:        bizID,
		ConversationID: conversationID,
		UserID:         strconv.FormatInt(userID, 10),
		ConnectorID:    connectorID,
		SectionID:      sectionID,
	})
	if err != nil {
		return nil, err
	}

	roundID := runRecord.ID

	userMessage, err := toConversationMessage(ctx, bizID, conversationID, userID, roundID, sectionID, message.MessageTypeQuestion, lastUserMessage)
	if err != nil {
		return nil, err
	}

	messageClient := crossmessage.DefaultSVC()
	_, err = messageClient.Create(ctx, userMessage)
	if err != nil {
		return nil, err
	}

	info, existed, unbinding, err := GetWorkflowDomainSVC().GetConvRelatedInfo(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	userSchemaMessage, err := toSchemaMessage(ctx, lastUserMessage)
	if err != nil {
		return nil, err
	}

	if existed {
		sr, err := GetWorkflowDomainSVC().StreamResume(ctx, &entity.ResumeRequest{
			EventID:    info.EventID,
			ExecuteID:  info.ExecID,
			ResumeData: lastUserMessage.Content,
		}, workflowModel.ExecuteConfig{
			Operator:     userID,
			Mode:         ternary.IFElse(isDebug, workflowModel.ExecuteModeDebug, workflowModel.ExecuteModeRelease),
			ConnectorID:  connectorID,
			ConnectorUID: strconv.FormatInt(userID, 10),
			BizType:      workflowModel.BizTypeWorkflow,
		})

		if err != nil {
			unErr := unbinding()
			if unErr != nil {
				logs.CtxErrorf(ctx, "unbinding failed, error: %v", unErr)
			}
			return nil, err
		}
		return schema.StreamReaderWithConvert(sr, w.convertToChatFlowRunResponseList(ctx, convertToChatFlowInfo{
			bizID:            bizID,
			conversationID:   conversationID,
			roundID:          roundID,
			workflowID:       workflowID,
			sectionID:        sectionID,
			unbinding:        unbinding,
			userMessage:      userSchemaMessage,
			suggestReplyInfo: req.GetSuggestReplyInfo(),
		})), nil

	}

	exeCfg := workflowModel.ExecuteConfig{
		ID:            mustParseInt64(req.GetWorkflowID()),
		From:          locator,
		Version:       version,
		Operator:      userID,
		Mode:          ternary.IFElse(isDebug, workflowModel.ExecuteModeDebug, workflowModel.ExecuteModeRelease),
		AppID:         appID,
		AgentID:       agentID,
		ConnectorID:   connectorID,
		ConnectorUID:  strconv.FormatInt(userID, 10),
		TaskType:      workflowModel.TaskTypeForeground,
		SyncPattern:   workflowModel.SyncPatternStream,
		InputFailFast: true,
		BizType:       workflowModel.BizTypeWorkflow,

		ConversationID: ptr.Of(conversationID),
		RoundID:        ptr.Of(roundID),
		InitRoundID:    ptr.Of(roundID),
		SectionID:      ptr.Of(sectionID),

		UserMessage: userSchemaMessage,
		Cancellable: isDebug,
	}

	historyMessages, err := makeChatFlowHistoryMessages(ctx, bizID, conversationID, userID, sectionID, connectorID, messages[:len(req.GetAdditionalMessages())-1])
	if err != nil {
		return nil, err
	}

	if len(historyMessages) > 0 {
		g := taskgroup.NewTaskGroup(ctx, len(historyMessages))
		for _, hm := range historyMessages {
			hMsg := hm
			g.Go(func() error {
				_, err := messageClient.Create(ctx, hMsg)
				if err != nil {
					return err
				}
				return nil
			})
		}
		err = g.Wait()
		if err != nil {
			logs.CtxWarnf(ctx, "create history message failed, err=%v", err)
		}
	}
	parameters[vo.UserInputKey], err = w.makeChatFlowUserInput(ctx, lastUserMessage)
	if err != nil {
		return nil, err
	}

	sr, err := GetWorkflowDomainSVC().StreamExecute(ctx, exeCfg, parameters)
	if err != nil {
		return nil, err
	}

	return schema.StreamReaderWithConvert(sr, w.convertToChatFlowRunResponseList(ctx, convertToChatFlowInfo{
		bizID:            bizID,
		conversationID:   conversationID,
		roundID:          roundID,
		workflowID:       workflowID,
		sectionID:        sectionID,
		unbinding:        unbinding,
		userMessage:      userSchemaMessage,
		suggestReplyInfo: req.GetSuggestReplyInfo(),
	})), nil

}
