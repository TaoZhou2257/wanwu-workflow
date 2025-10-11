package workflow

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/config"
	crossmodelmgrImpl "github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/service/crossdomain/impl/modelmgr"
	crosssearchImpl "github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/service/crossdomain/impl/search"
	crossuserImpl "github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/service/crossdomain/impl/user"
	"github.com/UnicomAI/wanwu/pkg/log"
	coze_app_modelmgr "github.com/coze-dev/coze-studio/backend/application/modelmgr"
	coze_app_upload "github.com/coze-dev/coze-studio/backend/application/upload"
	coze_app_user "github.com/coze-dev/coze-studio/backend/application/user"
	coze_app_workflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	coze_cross_user "github.com/coze-dev/coze-studio/backend/crossdomain/contract/user"
	coze_workflow "github.com/coze-dev/coze-studio/backend/domain/workflow"
	coze_workflow_service "github.com/coze-dev/coze-studio/backend/domain/workflow/service"
	coze_cache "github.com/coze-dev/coze-studio/backend/infra/cache/impl/redis"
	coze_checkpoint "github.com/coze-dev/coze-studio/backend/infra/checkpoint"
	coze_code "github.com/coze-dev/coze-studio/backend/infra/coderunner"
	coze_code_impl "github.com/coze-dev/coze-studio/backend/infra/coderunner/impl"
	coze_idgen "github.com/coze-dev/coze-studio/backend/infra/idgen/impl/idgen"
	coze_imagex "github.com/coze-dev/coze-studio/backend/infra/imagex"
	coze_storage "github.com/coze-dev/coze-studio/backend/infra/storage"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	workflowRepoSqlFile = "configs/schema.sql"
)

var _workflowService coze_workflow.Service

type Infra struct {
	DB      *gorm.DB
	Cache   *redis.Client
	Storage coze_storage.Storage
	ImageX  coze_imagex.ImageX
}

func Init(ctx context.Context, infra Infra) error {
	if _workflowService != nil {
		return errors.New("already init")
	}

	// init repo data
	if err := initRepo(infra.DB); err != nil {
		return fmt.Errorf("init repo err: %v", err)
	}

	// register all node adaptors
	coze_workflow_service.RegisterAllNodeAdaptors()
	coze_workflow_service.RegisterWanwuAllNodeAdaptors()

	// infra cache
	cache := coze_cache.NewWithRedisCli(infra.Cache)
	// id generator
	idGen, _ := coze_idgen.New(cache)
	// check point store
	cps := coze_checkpoint.NewRedisStore(cache)
	// code runner
	coze_code.SetCodeRunner(coze_code_impl.New())
	// workflow repo
	workflowRepo, err := coze_workflow_service.NewWorkflowRepository(idGen, infra.DB, cache, infra.Storage, cps, nil, config.Cfg().Workflow)
	if err != nil {
		return fmt.Errorf("init workflow err: %v", err)
	}
	coze_workflow.SetRepository(workflowRepo)

	// domain workflow service
	_workflowService = coze_workflow_service.NewWorkflowService(workflowRepo)

	// init application upload
	coze_app_upload.InitService(&coze_app_upload.UploadComponents{Cache: cache, Oss: infra.Storage, DB: infra.DB, Idgen: idGen})
	// init application modelmgr
	_ = coze_app_modelmgr.InitService(crossmodelmgrImpl.DefaultMock(), nil)
	// init application user
	_ = coze_app_user.InitService(ctx, infra.DB, infra.Storage, idGen)
	// init application workflow
	coze_app_workflow.SVC.DomainSVC = _workflowService
	coze_app_workflow.SVC.ImageX = infra.ImageX
	coze_app_workflow.SVC.TosClient = infra.Storage
	coze_app_workflow.SVC.IDGenerator = idGen
	coze_app_workflow.SetEventBus(crosssearchImpl.DefaultResourceEventBusMock())

	// init cross domain user
	coze_cross_user.SetDefaultSVC(crossuserImpl.DefaultMock())

	return nil
}

func initRepo(db *gorm.DB) error {
	// 读取sql文件
	file, err := os.Open(workflowRepoSqlFile)
	if err != nil {
		return fmt.Errorf("open %v err: %v", workflowRepoSqlFile, err)
	}
	defer file.Close()
	// 过滤注释行
	var b []byte
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("read %v err: %v", workflowRepoSqlFile, err)
		}

		lineStr := strings.TrimSpace(string(line))
		if strings.HasPrefix(lineStr, "-- ") {
			// 跳过注释行
			continue
		}
		b = append(b, line...)
		b = append(b, byte('\n'))
	}
	// 分隔sql，执行
	sqls := strings.Split(string(b), ";\n")
	for _, sql := range sqls {
		if sql = strings.TrimSpace(sql); sql == "" {
			continue
		}
		if err := db.Exec(sql).Error; err != nil {
			log.Warnf("execute [%v] err: %v", sql, err)
		}
	}
	return nil
}
