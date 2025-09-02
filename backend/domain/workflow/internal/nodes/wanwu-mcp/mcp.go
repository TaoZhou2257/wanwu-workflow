package wanwu_mcp

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/canvas/convert"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
)

const (
	WanWuMCPResult = "result"
)

type Config struct {
	SseUrl      string
	McpToolName string
}

func (c *Config) Adapt(_ context.Context, n *vo.Node, _ ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	inputs := n.Data.Inputs
	mcpToolInfoList := inputs.McpToolInfoList

	if len(mcpToolInfoList) > 0 {
		mcpInfo := mcpToolInfoList[0]
		c.SseUrl = mcpInfo.MCPServerURL
		c.McpToolName = mcpInfo.ToolName
	} else {
		return nil, errors.New("未找到MCP工具配置信息")
	}

	ns := &schema.NodeSchema{
		Key:     vo.NodeKey(n.ID),
		Type:    entity.NodeTypeWanWuMCPTool,
		Name:    n.Data.Meta.Title,
		Configs: c,
	}

	if err := convert.SetInputsForNodeSchema(n, ns); err != nil {
		return nil, err
	}

	if err := convert.SetOutputTypesForNodeSchema(n, ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (c *Config) Build(_ context.Context, ns *schema.NodeSchema, _ ...schema.BuildOption) (any, error) {
	if c.SseUrl == "" {
		return nil, errors.New("sse url is required")
	}

	if c.McpToolName == "" {
		return nil, errors.New("mcp tool name is required")
	}

	tool := &WanWuMCPTool{
		sseUrl:      c.SseUrl,
		mcpToolName: c.McpToolName,
	}
	return tool, nil
}

type WanWuMCPTool struct {
	sseUrl      string
	mcpToolName string
	mcpToolArgs map[string]any
}

func (i *WanWuMCPTool) Invoke(ctx context.Context, in map[string]any) (map[string]any, error) {
	i.mcpToolArgs = in

	transportClient, err := transport.NewSSEClientTransport(i.sseUrl,
		transport.WithSSEClientOptionReceiveTimeout(time.Minute*2),
		transport.WithSSEClientOptionLogger(logs.DefaultLogger()),
		transport.WithSSEClientOptionHTTPClient(&http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}))
	if err != nil {
		return nil, err
	}

	mcpClient, err := client.NewClient(transportClient)
	if err != nil {
		return nil, err
	}
	defer mcpClient.Close()

	request := &protocol.CallToolRequest{
		Name:      i.mcpToolName,
		Arguments: i.mcpToolArgs,
	}

	toolCallResult, err := mcpClient.CallTool(ctx, request)
	if err != nil {
		return nil, err
	}

	resultBytes, err := json.Marshal(toolCallResult)
	if err != nil {
		return nil, err
	}

	var resultMap map[string]any
	if err := json.Unmarshal(resultBytes, &resultMap); err != nil {
		return nil, err
	}

	return map[string]any{
		"result": resultMap,
	}, nil
}
