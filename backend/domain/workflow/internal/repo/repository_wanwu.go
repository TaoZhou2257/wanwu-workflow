package repo

import (
	einoCompose "github.com/cloudwego/eino/compose"
	"github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/repo/dal/query"
	"github.com/coze-dev/coze-studio/backend/infra/cache"
	cm "github.com/coze-dev/coze-studio/backend/infra/chatmodel"
	"github.com/coze-dev/coze-studio/backend/infra/idgen"
	"github.com/coze-dev/coze-studio/backend/infra/storage"
	"gorm.io/gorm"
)

// NewRepositoryWanwu 参考NewRepository，去掉chatModel初始化Suggester
func NewRepositoryWanwu(idgen idgen.IDGenerator, db *gorm.DB, redis cache.Cmdable, tos storage.Storage,
	cpStore einoCompose.CheckPointStore, chatModel cm.BaseChatModel, workflowConfig workflow.WorkflowConfig) (workflow.Repository, error) {
	return &RepositoryImpl{
		IDGenerator:     idgen,
		query:           query.Use(db),
		redis:           redis,
		tos:             tos,
		CheckPointStore: cpStore,
		InterruptEventStore: &interruptEventStoreImpl{
			redis: redis,
		},
		CancelSignalStore: &cancelSignalStoreImpl{
			redis: redis,
		},
		ExecuteHistoryStore: &executeHistoryStoreImpl{
			query: query.Use(db),
			redis: redis,
		},

		builtinModel:   chatModel,
		Suggester:      nil,
		WorkflowConfig: workflowConfig,
	}, nil

}
