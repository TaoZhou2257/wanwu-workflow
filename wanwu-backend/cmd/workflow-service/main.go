package main

import (
	"context"
	"flag"
	"fmt"
	"runtime"

	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/config"
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/server/http"
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/service/workflow"
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/coze-dev/coze-studio/backend/infra/impl/storage/minio"
)

var (
	configFile string
	isVersion  bool

	buildTime    string //编译时间
	buildVersion string //编译版本
	gitCommitID  string //git的commit id
	gitBranch    string //git branch
	builder      string //构建者
)

func main() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "conf yaml file")
	flag.BoolVar(&isVersion, "v", false, "build message")
	flag.Parse()

	if isVersion {
		versionPrint()
		return
	}

	ctx := context.Background()

	flag.Parse()
	if err := config.LoadConfig(configFile); err != nil {
		log.Fatalf("init cfg err: %v", err)
	}

	if err := log.InitLog(config.Cfg().Log.Std, config.Cfg().Log.Level, config.Cfg().Log.Logs...); err != nil {
		log.Fatalf("init log err: %v", err)
	}

	// db
	dbCli, err := db.New(config.Cfg().DB)
	if err != nil {
		log.Fatalf("init db err: %v", err)
	}

	// redis
	if err := redis.InitWorkflow(ctx, config.Cfg().Redis); err != nil {
		log.Fatalf("init redis err: %v", err)
	}

	// minio
	minioCli, err := minio.New(ctx, config.Cfg().Minio.Endpoint, config.Cfg().Minio.User, config.Cfg().Minio.Password, config.WorkflowBucket, false)
	if err != nil {
		log.Fatalf("init minio err: %v", err)
	}

	// workflow
	if err := workflow.Init(workflow.Config{DB: dbCli, Cache: redis.Workflow(), Storage: minioCli}); err != nil {
		log.Fatalf("init workflow service err: %v", err)
	}

	// http
	http.Init()
}

func versionPrint() {
	fmt.Printf("build_time: %s\n", buildTime)
	fmt.Printf("build_version: %s\n", buildVersion)
	fmt.Printf("git_commit_id: %s\n", gitCommitID)
	fmt.Printf("git branch: %s\n", gitBranch)
	fmt.Printf("runtime version: %s\n", runtime.Version())
	fmt.Printf("builder: %s\n", builder)
}
