package vo

type WanWuMCPTool struct {
	McpToolInfoList []*WanwuMCPToolInfo `json:"mcpInfoList"`
}

type WanwuMCPToolInfo struct {
	MCPServerURL string `json:"serverUrl"`
	ToolName     string `json:"name"`
}
