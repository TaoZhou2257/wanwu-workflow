package config

import (
	"github.com/UnicomAI/wanwu-workflow/wanwu-backend/pkg/redis"
	"github.com/UnicomAI/wanwu/pkg/db"
	"github.com/UnicomAI/wanwu/pkg/log"
	"github.com/UnicomAI/wanwu/pkg/util"
)

var (
	_c *Config
)

type Config struct {
	Server ServerConfig `json:"server" mapstructure:"server"`
	Log    LogConfig    `json:"log" mapstructure:"log"`
	DB     db.Config    `json:"db" mapstructure:"db"`
	Redis  redis.Config `json:"redis" mapstructure:"redis"`
	Minio  MinioConfig  `json:"minio" mapstructure:"minio"`
}

type ServerConfig struct {
	HttpEndpoint   string `json:"http_endpoint" mapstructure:"http_endpoint"`
	MaxReqBodySize int    `json:"max_req_body_size" mapstructure:"max_req_body_size"`
}

type LogConfig struct {
	Std   bool         `json:"std" mapstructure:"std"`
	Level string       `json:"level" mapstructure:"level"`
	Logs  []log.Config `json:"logs" mapstructure:"logs"`
}

type MinioConfig struct {
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
	User     string `json:"user" mapstructure:"user"`
	Password string `json:"password" mapstructure:"password"`
}

func LoadConfig(in string) error {
	_c = &Config{}
	return util.LoadConfig(in, _c)
}

func Cfg() *Config {
	if _c == nil {
		log.Panicf("cfg nil")
	}
	return _c
}
