package adaptor

import (
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/filegenerator"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/fileparser"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/knowledge"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/mcp"
	wanwu_intentdetector "github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes/wanwu-intentdetector"
)

// RegisterWanwuAllNodeAdaptors 参考RegisterAllNodeAdaptors
func RegisterWanwuAllNodeAdaptors() {
	// register a generator function so that each time a NodeAdaptor is needed,
	// we can provide a brand new Config instance.
	nodes.RegisterNodeAdaptor(entity.NodeTypeWanwuIntentDetector, func() nodes.NodeAdaptor {
		return &wanwu_intentdetector.Config{}
	})

	nodes.RegisterNodeAdaptor(entity.NodeTypeWanWuKnowledgeRetriever, func() nodes.NodeAdaptor {
		return &knowledge.WanWuRetrieveConfig{}
	})

	nodes.RegisterNodeAdaptor(entity.NodeTypeWanWuFileParser, func() nodes.NodeAdaptor {
		return &fileparser.WanWuRetrieveConfig{}
	})

	nodes.RegisterNodeAdaptor(entity.NodeTypeWanWuFileGenerator, func() nodes.NodeAdaptor {
		return &filegenerator.WanWuRetrieveConfig{}
	})

	nodes.RegisterNodeAdaptor(entity.NodeTypeWanWuMCPTool, func() nodes.NodeAdaptor {
		return &mcp.Config{}
	})

	// register branch adaptors
	nodes.RegisterBranchAdaptor(entity.NodeTypeWanwuIntentDetector, func() nodes.BranchAdaptor {
		return &wanwu_intentdetector.Config{}
	})
}
