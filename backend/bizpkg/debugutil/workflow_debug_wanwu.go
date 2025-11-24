package debugutil

import (
	"context"
	"net/url"
	"os"
	"strconv"

	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

func GetWorkflowDebugURLByWanwu(ctx context.Context, workflowID, spaceID, executeID int64) string {
	workFlowURL, err := url.JoinPath(os.Getenv("WANWU_EXTERNAL_SCHEME")+"://"+os.Getenv("WANWU_EXTERNAL_ENDPOINT"), "work_flow")
	if err != nil {
		logs.CtxErrorf(ctx, "[GetWorkflowDebugURLByWanwu] join path failed, use default debug url instead, err: %v", err)
		return ""
	}

	u, err := url.Parse(workFlowURL)
	if err != nil {
		logs.CtxErrorf(ctx, "[GetWorkflowDebugURLByWanwu] parse workflow url failed, use default debug url instead, err: %v", err)
		return ""
	}

	q := u.Query()
	q.Set("execute_id", strconv.FormatInt(executeID, 10))
	q.Set("space_id", strconv.FormatInt(spaceID, 10))
	q.Set("workflow_id", strconv.FormatInt(workflowID, 10))
	q.Set("execute_mode", "2")
	u.RawQuery = q.Encode()

	return u.String()
}
