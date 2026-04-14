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
	WanWuMCPResult         = "result"
	MCPTransportSSE        = "sse"
	MCPTransportStreamable = "streamable"
)

type Config struct {
	SseUrl        string
	McpToolName   string
	Transport     string // 传输协议: "sse" 或 "streamable"
	StreamableUrl string // Streamable HTTP URL
}

func (c *Config) Adapt(_ context.Context, n *vo.Node, _ ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	inputs := n.Data.Inputs
	mcpToolInfoList := inputs.McpToolInfoList

	if len(mcpToolInfoList) > 0 {
		mcpInfo := mcpToolInfoList[0]
		c.SseUrl = mcpInfo.MCPServerURL
		c.McpToolName = mcpInfo.ToolName
		c.Transport = mcpInfo.Transport
		c.StreamableUrl = mcpInfo.StreamableURL
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
	// 根据 transport 类型选择 URL
	var serverUrl string
	var transportType string
	switch c.Transport {
	case MCPTransportStreamable:
		serverUrl = c.StreamableUrl
		transportType = MCPTransportStreamable
	case MCPTransportSSE:
		serverUrl = c.SseUrl
		transportType = MCPTransportSSE
	default:
		return nil, errors.New("transport not support")
	}

	if serverUrl == "" {
		return nil, errors.New("server url is required")
	}

	if c.McpToolName == "" {
		return nil, errors.New("mcp tool name is required")
	}

	tool := &WanWuMCPTool{
		serverUrl:     serverUrl,
		mcpToolName:   c.McpToolName,
		transportType: transportType,
	}
	return tool, nil
}

type WanWuMCPTool struct {
	serverUrl     string
	mcpToolName   string
	transportType string
	mcpToolArgs   map[string]any
}

// httpClient 创建共享的 HTTP 客户端，跳过证书验证
var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	},
}

func (i *WanWuMCPTool) Invoke(ctx context.Context, in map[string]any) (map[string]any, error) {
	i.mcpToolArgs = in

	var transportClient transport.ClientTransport
	var err error

	switch i.transportType {
	case MCPTransportStreamable:
		// 创建 StreamableHTTP 传输客户端
		transportClient, err = transport.NewStreamableHTTPClientTransport(i.serverUrl,
			transport.WithStreamableHTTPClientOptionLogger(logs.DefaultLogger()),
			transport.WithStreamableHTTPClientOptionHTTPClient(httpClient),
		)
		if err != nil {
			return nil, err
		}
	case MCPTransportSSE:
		// 默认使用 SSE 传输客户端
		transportClient, err = transport.NewSSEClientTransport(i.serverUrl,
			transport.WithSSEClientOptionReceiveTimeout(time.Minute*2),
			transport.WithSSEClientOptionLogger(logs.DefaultLogger()),
			transport.WithSSEClientOptionHTTPClient(httpClient),
		)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("transport not support")
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
