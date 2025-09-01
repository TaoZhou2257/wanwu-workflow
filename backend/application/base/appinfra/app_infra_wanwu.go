package appinfra

import "github.com/coze-dev/coze-studio/backend/infra/contract/coderunner"

func InitCodeRunner() coderunner.Runner {
	return initCodeRunner()
}
