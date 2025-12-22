package repo

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
)

type executeHistoryStoreImplByWanwu struct {
	*executeHistoryStoreImpl
}

func (e *executeHistoryStoreImplByWanwu) CreateWorkflowExecution(ctx context.Context, execution *entity.WorkflowExecution) (err error) {
	if skipRecordWorkflowExecuteHistory(ctx) {
		return nil
	}
	return e.executeHistoryStoreImpl.CreateWorkflowExecution(ctx, execution)
}

func (e *executeHistoryStoreImplByWanwu) UpdateWorkflowExecution(ctx context.Context, execution *entity.WorkflowExecution,
	allowedStatus []entity.WorkflowExecuteStatus) (_ int64, _ entity.WorkflowExecuteStatus, err error) {
	if skipRecordWorkflowExecuteHistory(ctx) {
		return 1, entity.WorkflowSuccess, nil
	}
	return e.executeHistoryStoreImpl.UpdateWorkflowExecution(ctx, execution, allowedStatus)
}

func (e *executeHistoryStoreImplByWanwu) CreateNodeExecution(ctx context.Context, execution *entity.NodeExecution) error {
	if skipRecordWorkflowExecuteHistory(ctx) {
		return nil
	}
	return e.executeHistoryStoreImpl.CreateNodeExecution(ctx, execution)
}

func (e *executeHistoryStoreImplByWanwu) UpdateNodeExecution(ctx context.Context, execution *entity.NodeExecution) (err error) {
	if skipRecordWorkflowExecuteHistory(ctx) {
		return nil
	}
	return e.executeHistoryStoreImpl.UpdateNodeExecution(ctx, execution)
}

func skipRecordWorkflowExecuteHistory(ctx context.Context) bool {
	val := ctx.Value("WANWU_WORKFLOW_OPENAPI_RUN_SKIP_EXECUTE_HISTORY")
	if val == nil {
		return false
	}
	skip, ok := val.(bool)
	if !ok {
		return false
	}
	return skip
}
