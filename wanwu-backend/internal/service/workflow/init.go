package workflow

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/config"
	crosssearchImpl "github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/service/crossdomain/impl/search"
	crossuserImpl "github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/service/crossdomain/impl/user"
	"github.com/UnicomAI/wanwu/pkg/log"
	coze_app_upload "github.com/coze-dev/coze-studio/backend/application/upload"
	coze_app_user "github.com/coze-dev/coze-studio/backend/application/user"
	coze_app_workflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	coze_cross_upload "github.com/coze-dev/coze-studio/backend/crossdomain/upload"
	coze_cross_upload_impl "github.com/coze-dev/coze-studio/backend/crossdomain/upload/impl"
	coze_cross_user "github.com/coze-dev/coze-studio/backend/crossdomain/user"
	coze_workflow "github.com/coze-dev/coze-studio/backend/domain/workflow"
	coze_workflow_service "github.com/coze-dev/coze-studio/backend/domain/workflow/service"
	coze_cache "github.com/coze-dev/coze-studio/backend/infra/cache"
	coze_checkpoint "github.com/coze-dev/coze-studio/backend/infra/checkpoint"
	coze_code "github.com/coze-dev/coze-studio/backend/infra/coderunner"
	coze_code_impl "github.com/coze-dev/coze-studio/backend/infra/coderunner/impl"
	coze_idgen "github.com/coze-dev/coze-studio/backend/infra/idgen/impl/idgen"
	coze_imagex "github.com/coze-dev/coze-studio/backend/infra/imagex"
	coze_storage "github.com/coze-dev/coze-studio/backend/infra/storage"
	"gorm.io/gorm"
)

const (
	workflowRepoSqlFile         = "configs/schema.sql"
	workflowRepoSqlMigrationDir = "configs/migrations"
)

var _workflowService coze_workflow.Service

type Infra struct {
	DB      *gorm.DB
	Cache   coze_cache.Cmdable
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
	// migrate repo
	if err := migrateRepo(infra.DB); err != nil {
		return fmt.Errorf("migrate repo err: %v", err)
	}

	// register all node adaptors
	coze_workflow_service.RegisterAllNodeAdaptors()
	coze_workflow_service.RegisterWanwuAllNodeAdaptors()

	// id generator
	idGen, _ := coze_idgen.New(infra.Cache)
	// check point store
	cps := coze_checkpoint.NewRedisStore(infra.Cache)
	// code runner
	coze_code.SetCodeRunner(coze_code_impl.NewByWanwu())
	// workflow repo
	workflowRepo, _ := coze_workflow_service.NewWorkflowRepositoryWanwu(idGen, infra.DB, infra.Cache, infra.Storage, cps, nil, config.Cfg().Workflow)
	coze_workflow.SetRepository(workflowRepo)

	// domain workflow service
	_workflowService = coze_workflow_service.NewWorkflowService(workflowRepo)

	// init application upload
	coze_app_upload.InitService(&coze_app_upload.UploadComponents{Cache: infra.Cache, Oss: infra.Storage, DB: infra.DB, Idgen: idGen})
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
	coze_cross_upload.SetDefaultWanwuSVC(coze_cross_upload_impl.NewWanwuUploader(infra.Storage))

	return nil
}

func initRepo(db *gorm.DB) error {
	return executeSqlFile(db, workflowRepoSqlFile, true)

}

func migrateRepo(db *gorm.DB) error {
	// 读取migrations中的sql文件
	if err := filepath.Walk(workflowRepoSqlMigrationDir, func(filePath string, fileInfo fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if fileInfo.IsDir() {
			return nil
		}
		return executeSqlFile(db, filePath, true)
	}); err != nil {
		return err
	}
	return nil
}

func executeSqlFile(db *gorm.DB, sqlFilePath string, skipExecErr bool) error {
	// 读取sql文件
	file, err := os.Open(sqlFilePath)
	if err != nil {
		return fmt.Errorf("open %v err: %v", sqlFilePath, err)
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
			return fmt.Errorf("read %v err: %v", sqlFilePath, err)
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
			if !skipExecErr {
				return fmt.Errorf("execute [%v] err: %v", sql, err)
			}
			log.Warnf("execute [%v] err: %v", sql, err)
		}
	}
	return nil
}
