package coze

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/coze-dev/coze-studio/backend/api/model/workflow"
	appworkflow "github.com/coze-dev/coze-studio/backend/application/workflow"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

// ImportWorkFlow .
// @router /api/workflow_api/import [POST]
func ImportWorkFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req importWorkflowRequest
	var flowMode *workflow.WorkflowMode
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	// create
	createReq := workflow.CreateWorkflowRequest{
		SpaceID:  req.SpaceID,
		Name:     req.Name,
		Desc:     req.Desc,
		FlowMode: flowMode,
	}
	switch req.FlowMode {
	case "3":
		createReq.FlowMode = workflow.WorkflowModePtr(workflow.WorkflowMode_ChatFlow)
	default:
		createReq.FlowMode = workflow.WorkflowModePtr(workflow.WorkflowMode_Workflow)
	}
	if req.IconUrl != "" {
		createReq.IconURI = req.IconUrl
	}
	resp, err := appworkflow.SVC.CreateWorkflowByWanwu(ctx, &createReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	// save
	saveReq := workflow.SaveWorkflowRequest{
		WorkflowID: resp.Data.WorkflowID,
		SpaceID:    &req.SpaceID,
		Schema:     &req.Schema,
	}
	_, err = appworkflow.SVC.SaveWorkflow(ctx, &saveReq)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	c.JSON(consts.StatusOK, resp)
}

type importWorkflowRequest struct {
	// Space id, cannot be empty
	SpaceID string `form:"space_id,required" json:"space_id" query:"space_id,required"`
	// process name
	Name string `form:"name,required" json:"name" query:"name,required"`
	// Process description, not null
	Desc string `form:"desc,required" json:"desc" query:"desc,required"`
	// file data
	Schema string `form:"schema" json:"schema" query:"schema"`
	// icon url
	IconUrl string `form:"icon_url" json:"icon_url" query:"icon_url"`
	// flow mode
	FlowMode string `form:"flow_mode" json:"flow_mode" query:"flow_mode"`
}

// ExportWorkFlow .
// @router /api/workflow_api/export [POST]
func ExportWorkFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req appworkflow.ExportWorkflowRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.GetCanvasInfoByWanwu(ctx, &req)
	if err != nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	if resp.Data.Workflow.SchemaJSON == nil {
		internalServerErrorResponse(ctx, c, err)
		return
	}
	// 解析原始schema
	var schema vo.Canvas
	if err := sonic.Unmarshal([]byte(*resp.Data.Workflow.SchemaJSON), &schema); err != nil {
		internalServerErrorResponse(ctx, c, fmt.Errorf("failed to parse schema JSON: %v", err))
		return
	}
	// 清理所有节点的用户自定义参数
	var schemaStr string
	if err := cleanWorkflowSchema(&schema); err != nil {
		internalServerErrorResponse(ctx, c, fmt.Errorf("failed to clean workflow schema: %v", err))
		return
	}
	// Wanwu 导出后处理：
	// - HTTP 节点 `authOpen=true`：清空 BEARER/CUSTOM/BASIC 的 token 值
	// - 并补齐每个鉴权数组元素的顶层 `type`，避免前端初始化时报 `param.type === undefined`
	//
	// 这里在“序列化后的 map 层”处理，避免改动共享的 VO 结构体。
	schemaBytes, err := sonic.Marshal(schema)
	if err != nil {
		internalServerErrorResponse(ctx, c, fmt.Errorf("failed to marshal workflow schema bytes: %v", err))
		return
	}
	var schemaMap map[string]any
	if err := sonic.Unmarshal(schemaBytes, &schemaMap); err != nil {
		internalServerErrorResponse(ctx, c, fmt.Errorf("failed to unmarshal workflow schema map: %v", err))
		return
	}
	cleanWANWUHTTPAuthInSchemaJSON(schemaMap)

	if schemaStr, err = sonic.MarshalString(schemaMap); err != nil {
		internalServerErrorResponse(ctx, c, fmt.Errorf("failed to marshal workflow schema string: %v", err))
		return
	}
	// 创建导出数据
	exportData := workflowExport{
		WorkflowName: resp.Data.Workflow.Name,
		WorkflowDesc: resp.Data.Workflow.Desc,
		Schema:       schemaStr,
	}
	response := map[string]any{
		"data":    exportData,
		"code":    0,
		"message": "",
	}
	c.JSON(consts.StatusOK, response)
}

// --- internal ---

type workflowExport struct {
	WorkflowName string `json:"name"`
	WorkflowDesc string `json:"desc"`
	Schema       string `json:"schema"`
}

// cleanWorkflowSchema 清理整个工作流schema
func cleanWorkflowSchema(schema *vo.Canvas) error {
	if schema == nil {
		return nil
	}
	// 遍历所有节点
	for _, node := range schema.Nodes {
		cleanNode(node)
	}
	return nil
}

func cleanNode(node *vo.Node) {
	if node == nil {
		return
	}
	switch node.Type {
	case "3": // 大模型节点
		cleanLLMNode(node)
	case "21": // 循环节点
		cleanLoopNode(node)
	case "22": // 意图识别节点
		cleanIntentNode(node)
	case "12", "42", "43", "44", "46": // 数据库节点
		cleanDatabaseNode(node)
	case "1006": // 知识库检索节点
		cleanKnowledgeNode(node)
	case "1009": // MCP节点
		cleanMCPNode(node)
	case "1010": // GUI智能体节点
		cleanGUINode(node)
	case "1004": // Tool节点
		cleanToolNode(node)
	case "1012": // 问答库检索节点
		cleanQANode(node)
	case "1013": // 智能体节点
		cleanAgentNode(node)
	}
}

// cleanWANWUHTTPAuthInSchemaJSON：对 wanwu 导出的 workflow JSON 做 HTTP 鉴权清理与兼容补齐。
// param.type 这个字段在你们当前 wanwu 导出用到的 Go 结构里并没有对应字段（vo.Param 里只有 name/input/left/right/variables，没有 type）。
// 所以如果你把 node *vo.Node 传给处理函数，你在 Go struct 里加不了 param.type，最终 marshal 回 JSON 时也就补不回来，前端还是会拿到 undefined。
// 实现方式：
// - 直接在 map[string]any 层处理，避免修改共享的 VO 结构体。
func cleanWANWUHTTPAuthInSchemaJSON(schemaMap map[string]any) {
	if schemaMap == nil {
		return
	}

	nodes, ok := schemaMap["nodes"].([]any)
	if !ok {
		return
	}

	for _, nodeAny := range nodes {
		node, ok := nodeAny.(map[string]any)
		if !ok {
			continue
		}
		nodeType, _ := node["type"].(string)
		if nodeType != "45" { // HTTP 请求
			continue
		}

		data, ok := node["data"].(map[string]any)
		if !ok {
			continue
		}
		inputs, ok := data["inputs"].(map[string]any)
		if !ok {
			continue
		}
		auth, ok := inputs["auth"].(map[string]any)
		if !ok {
			continue
		}

		authOpen, _ := auth["authOpen"].(bool)
		if !authOpen {
			continue
		}

		authData, ok := auth["authData"].(map[string]any)
		if !ok {
			continue
		}

		// 清空单个鉴权条目的 token/值（并补齐前端需要的 param.type）。
		clearTokenValue := func(param map[string]any) {
			if param == nil {
				return
			}
			input, ok := param["input"].(map[string]any)
			if !ok {
				return
			}

			// 前端初始化需要顶层 param.type；wanwu 导出源只在 input.type 里提供。
			if inputType, ok := input["type"]; ok {
				param["type"] = inputType
			}

			value, ok := input["value"].(map[string]any)
			if !ok {
				return
			}
			value["content"] = ""
		}

		processParamList := func(container map[string]any, arrayKey string) {
			arr, ok := container[arrayKey].([]any)
			if !ok {
				return
			}
			for _, itemAny := range arr {
				item, ok := itemAny.(map[string]any)
				if !ok {
					continue
				}
				clearTokenValue(item)
			}
		}

		// 前端会遍历以下数组（只要导出 JSON 里存在相应字段）。
		processParamList(authData, "bearerTokenData")
		processParamList(authData, "basicAuthData")
		if customData, ok := authData["customData"].(map[string]any); ok {
			processParamList(customData, "data")
		}

	}
}

// cleanLLMNode 清理大模型节点 - 删除modelType和modelName
func cleanLLMNode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}

	// 直接操作 LLMParam
	if paramSlice, ok := node.Data.Inputs.LLMParam.([]interface{}); ok {
		var newParams []interface{}
		for _, item := range paramSlice {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if name, ok := itemMap["name"].(string); ok {
					if name == "modleName" || name == "modelType" {
						continue
					}
				}
			}
			newParams = append(newParams, item)
		}
		node.Data.Inputs.LLMParam = newParams
	}
}

// cleanLoopNode 清理循环节点
func cleanLoopNode(node *vo.Node) {
	if len(node.Blocks) == 0 {
		return
	}
	for _, node := range node.Blocks {
		cleanNode(node)
	}
}

// cleanIntentNode 清理意图识别节点 - 将modelType和modelName置为空
func cleanIntentNode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}

	// 适配对象形式的 llmParam
	if paramMap, ok := node.Data.Inputs.LLMParam.(map[string]interface{}); ok {
		// 创建新的 map，排除不需要的字段
		newParams := make(map[string]interface{})
		for key, value := range paramMap {
			if key != "modelType" && key != "modelName" {
				newParams[key] = value
			}
		}
		node.Data.Inputs.LLMParam = newParams
	}
}

// cleanKnowledgeNode 清理知识库检索节点 - 将knowledgeList置为空
func cleanKnowledgeNode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}

	// 遍历 InputParameters 找到 knowledgeList 并清空其 content
	for i := range node.Data.Inputs.DatasetParam {
		param := node.Data.Inputs.DatasetParam[i]
		if param.Name == "knowledgeList" {
			// 清空 content 切片
			if param.Input != nil && param.Input.Value != nil {
				param.Input.Value.Content = []interface{}{}
			}
			break
		}
	}
}

// cleanMCPNode 清理MCP节点
func cleanMCPNode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}

	// 检查是否有 WanWuMCPTool 配置
	if node.Data.Inputs.WanWuMCPTool != nil {
		// 将 mcpToolInfoList 置为空切片
		node.Data.Inputs.WanWuMCPTool.McpToolInfoList = make([]*vo.WanWuMCPToolInfo, 0)
	}
}

// cleanGUINode 清理GUI智能体节点 - 将modelId置为空
func cleanGUINode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}

	// 检查是否有 WanwuGUIParam 配置
	if node.Data.Inputs.WanwuGUIParam != nil {
		// 将 modelId 置为空
		node.Data.Inputs.WanwuGUIParam.ModelID = ""
	}
}

// cleanToolNode 清理Tool节点 - 删除apiKey和header-Authorization
func cleanToolNode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}
	// 清理 apiParam 中的 apiKey
	if node.Data.Inputs.APIParams != nil {
		for _, param := range node.Data.Inputs.APIParams {
			if param.Name == "apiKey" {
				// 删除该参数
				param.Input.Value.Content = ""
			}
		}
		for _, param := range node.Data.Inputs.InputParameters {
			if param.Name == "header-Authorization" {
				//删除该参数
				param.Input.Value.Content = ""
			}
		}
	}
	// 清理 inputParameters 中的敏感信息
	if node.Data.Inputs.InputParameters != nil {
		for _, param := range node.Data.Inputs.InputParameters {
			// 清理 query-key
			if param.Name == "query-key" {
				param.Input.Value.Content = ""
			}
		}
	}
	if node.Data.Inputs.WanwuToolParam != nil {
		node.Data.Inputs.WanwuToolParam.ApiKey = ""
	}
}

// cleanDatabaseNode 清理数据库节点 - 删除databaseInfoList
func cleanDatabaseNode(node *vo.Node) {
	if node == nil || node.Data == nil || node.Data.Inputs == nil {
		return
	}
	// 清空databaseInfoList字段
	node.Data.Inputs.DatabaseInfoList = nil
}

// cleanQANode 清理问答库检索节点（类型1012）- 只置空knowledgeList
func cleanQANode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}
	for _, param := range node.Data.Inputs.DatasetParam {
		if param != nil && param.Name == "knowledgeList" {
			// 置空knowledgeList的Input内容
			if param.Input != nil && param.Input.Value != nil {
				// 将Value的内容置为空数组
				param.Input.Value.Content = make([]any, 0)
			}
			break
		}
	}
}

// cleanAgentNode 清理智能体节点（类型1013）- 置空modelType和agentToolParams
func cleanAgentNode(node *vo.Node) {
	if node.Data == nil || node.Data.Inputs == nil {
		return
	}

	// 1. 将llmParam中的modelType值置空
	if paramSlice, ok := node.Data.Inputs.LLMParam.([]interface{}); ok {
		for _, item := range paramSlice {
			if itemMap, ok := item.(map[string]interface{}); ok {
				if name, ok := itemMap["name"].(string); ok && name == "modelType" {
					// 将modelType的值置空
					if input, ok := itemMap["input"].(map[string]interface{}); ok {
						if value, ok := input["value"].(map[string]interface{}); ok {
							value["content"] = ""
						}
					}
				}
			}
		}
	}

	// 2. 将agentToolParams置为空数组
	node.Data.Inputs.AgentToolParams = nil
}
