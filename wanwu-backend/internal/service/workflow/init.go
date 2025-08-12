package workflow

import (
	"errors"

	"github.com/coze-dev/coze-studio/backend/domain/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/service"
	"github.com/coze-dev/coze-studio/backend/infra/contract/cache"
	"github.com/coze-dev/coze-studio/backend/infra/contract/storage"
	"github.com/coze-dev/coze-studio/backend/infra/impl/checkpoint"
	"github.com/coze-dev/coze-studio/backend/infra/impl/idgen"
	"gorm.io/gorm"
)

var _workflowService workflow.Service

type Config struct {
	DB      *gorm.DB
	Cache   cache.Cmdable
	Storage storage.Storage
}

func Init(cfg Config) error {
	if _workflowService != nil {
		return errors.New("already init")
	}

	// id generator
	idGen, _ := idgen.New(cfg.Cache)
	// check point store
	cps := checkpoint.NewRedisStore(cfg.Cache)
	// workflow repo
	workflowRepo := service.NewWorkflowRepository(idGen, cfg.DB, cfg.Cache, cfg.Storage, cps, nil)
	workflow.SetRepository(workflowRepo)

	// workflow service
	_workflowService = service.NewWorkflowService(workflowRepo)
	return nil
}

func Service() workflow.Service {
	return _workflowService
}
