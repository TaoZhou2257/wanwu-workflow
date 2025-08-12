package http

import (
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/config"
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/server/http/router"
	hertz_server "github.com/cloudwego/hertz/pkg/app/server"
	hertz_config "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/coze-dev/coze-studio/backend/api/middleware"
	hertz_cors "github.com/hertz-contrib/cors"
)

func Init() {

	opts := []hertz_config.Option{
		hertz_server.WithHostPorts(config.Cfg().Server.HttpEndpoint),
		hertz_server.WithMaxRequestBodySize(config.Cfg().Server.MaxReqBodySize),
	}

	s := hertz_server.Default(opts...)

	// cors option
	corsCfg := hertz_cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = []string{"*"}
	corsHandler := hertz_cors.New(corsCfg)

	// Middleware order matters
	s.Use(middleware.ContextCacheMW())     // must be first
	s.Use(middleware.RequestInspectorMW()) // must be second
	s.Use(middleware.SetHostMW())
	s.Use(middleware.SetLogIDMW())
	s.Use(corsHandler)
	s.Use(middleware.AccessLogMW())
	s.Use(middleware.OpenapiAuthMW())
	s.Use(middleware.SessionAuthMW())
	s.Use(middleware.I18nMW()) // must after SessionAuthMW

	router.Register(s)
	s.Spin()
}
