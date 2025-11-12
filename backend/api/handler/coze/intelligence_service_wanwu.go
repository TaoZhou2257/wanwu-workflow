package coze

import (
	"context"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/app/intelligence"
	"github.com/coze-dev/coze-studio/backend/api/model/app/intelligence/common"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
)

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

	// 从header的Referer中，获取当前workflow id，固定为当前workflow对应的唯一intelligence的id
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

	// 根据workflow id，获取workflow信息，用于设置当前workflow对应的唯一ntelligence相关信息
	workflowIDStr := strconv.Itoa(workflowID)
	resp, err := appworkflow.SVC.GetCanvasInfo(ctx, &workflow.GetCanvasInfoRequest{
		SpaceID:    strconv.Itoa(int(req.SpaceID)),
		WorkflowID: &workflowIDStr,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	if resp.Data.Workflow.URL == "" || strings.Contains(resp.Data.Workflow.URL, "default_workflow_icon.png") {
		// 设置默认图标
		resp.Data.Workflow.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
			os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
	}

	c.JSON(consts.StatusOK, &intelligence.GetDraftIntelligenceListResponse{
		Data: &intelligence.DraftIntelligenceListData{
			Intelligences: []*intelligence.IntelligenceData{workflow2intelligenceData(resp.Data.Workflow)},
			Total:         1,
			HasMore:       false,
			NextCursorID:  "",
		},
	})
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

	// req.IntelligenceID 即 workflow id
	// 根据workflow id，获取workflow信息，用于设置当前workflow对应的唯一ntelligence相关信息
	workflowIDStr := strconv.Itoa(int(req.IntelligenceID))
	resp, err := appworkflow.SVC.GetCanvasInfo(ctx, &workflow.GetCanvasInfoRequest{
		SpaceID:    c.Request.Header.Get("X-Org-Id"),
		WorkflowID: &workflowIDStr,
	})
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	if resp.Data.Workflow.URL == "" || strings.Contains(resp.Data.Workflow.URL, "default_workflow_icon.png") {
		// 设置默认图标
		resp.Data.Workflow.URL, _ = url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"),
			os.Getenv("WANWU_WORKFLOW_DEFAULT_ICON"))
	}

	intellData := workflow2intelligenceData(resp.Data.Workflow)
	c.JSON(consts.StatusOK, &intelligence.GetDraftIntelligenceInfoResponse{
		Data: &intelligence.GetDraftIntelligenceInfoData{
			IntelligenceType: intellData.Type,
			BasicInfo:        intellData.BasicInfo,
			PublishInfo:      intellData.PublishInfo,
			OwnerInfo:        intellData.OwnerInfo,
		},
		Code:     0,
		Msg:      "",
		BaseResp: nil,
	})
}

func workflow2intelligenceData(workflow *workflow.Workflow) *intelligence.IntelligenceData {
	workflowID, _ := strconv.Atoi(workflow.WorkflowID)
	spaceID, _ := strconv.Atoi(*workflow.SpaceID)
	createID, _ := strconv.Atoi(workflow.Creator.ID)
	return &intelligence.IntelligenceData{
		BasicInfo: &common.IntelligenceBasicInfo{
			ID:          int64(workflowID),
			Name:        workflow.Name,
			Description: workflow.Desc,
			IconURI:     workflow.IconURI,
			IconURL:     workflow.URL,
			SpaceID:     int64(spaceID),
			OwnerID:     int64(createID),
			CreateTime:  workflow.CreateTime,
			UpdateTime:  workflow.UpdateTime,
			Status:      common.IntelligenceStatus_Using,
			PublishTime: 0,
		},
		Type: common.IntelligenceType_Project,
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
			UserID:         int64(createID),
			Nickname:       workflow.Creator.Name,
			AvatarURL:      workflow.Creator.AvatarURL,
			UserUniqueName: workflow.Creator.Name,
			UserLabel:      nil,
		},
		LatestAuditInfo: nil,
		FavoriteInfo: &intelligence.FavoriteInfo{
			IsFav:   false,
			FavTime: "",
		},
		OtherInfo: nil,
	}
}
