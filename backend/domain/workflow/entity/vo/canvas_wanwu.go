package vo

type WanWuMCPTool struct {
	McpToolInfoList []*WanWuMCPToolInfo `json:"mcpInfoList"`
}

type WanWuMCPToolInfo struct {
	MCPServerURL string `json:"serverUrl"`
	ToolName     string `json:"name"`
}

type WanWuGUIParam struct {
	ModelID string `json:"modelId"`
}
