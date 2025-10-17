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

type WanWuTool struct {
	ToolID     string `json:"toolId"`
	ApiKey     string `json:"apiKey"`
	ToolType   string `json:"toolType"` // 内置:"builtin" 自定义:"custom"
	ToolName   string `json:"toolName"`
	ActionName string `json:"actionName"`
	ActionID   string `json:"actionId"`
}
