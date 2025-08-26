package workflow

import "github.com/coze-dev/coze-studio/backend/domain/search/service"

func SetEventBus(bus service.ResourceEventBus) {
	eventBus = bus
}
