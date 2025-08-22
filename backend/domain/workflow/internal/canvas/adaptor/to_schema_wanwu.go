package adaptor

import (
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	wanwu_intentdetector "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/wanwu-intentdetector"
)

// RegisterWanwuAllNodeAdaptors 参考RegisterAllNodeAdaptors
func RegisterWanwuAllNodeAdaptors() {
	// register a generator function so that each time a NodeAdaptor is needed,
	// we can provide a brand new Config instance.
	nodes.RegisterNodeAdaptor(entity.NodeTypeWanwuIntentDetector, func() nodes.NodeAdaptor {
		return &wanwu_intentdetector.Config{}
	})

	// register branch adaptors
	nodes.RegisterBranchAdaptor(entity.NodeTypeWanwuIntentDetector, func() nodes.BranchAdaptor {
		return &wanwu_intentdetector.Config{}
	})
}
