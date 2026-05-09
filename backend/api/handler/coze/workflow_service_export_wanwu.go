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
	// Wanwu 导出后处理（HTTP 节点的特殊路径）：
	// - HTTP 节点 `authOpen=true`：清空 BEARER/CUSTOM/BASIC 的 token 值
	// - 并补齐每个鉴权数组元素的顶层 `type`，避免前端初始化时报 `param.type === undefined`
	//
	// 为什么不放进上面的 cleanNode switch（即 *vo.Node 那条链路）：
	//   1. 前端期望的 `param.type` 在 vo.Param 里根本没有对应字段，传 *vo.Node 进去
	//      即使设上也会在 marshal 时丢失，前端依旧会拿到 undefined。
	//   2. wanwu 实际导出 JSON 里 auth 结构（含 basicAuthData、input.value.content 等）
	//      和 vo.Auth 也并不完全对齐，先反序列化到结构体会丢字段。
	// 因此这里把 schema 先 marshal 成 JSON、再 unmarshal 成 map[string]any，
	// 在“序列化后的 map 层”做修改，既不污染共享 VO，又能保留/补齐所有原始字段。
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

// cleanWANWUHTTPAuthInSchemaJSON 对 wanwu 导出的 workflow JSON 做 HTTP 节点（type="45"）的鉴权清理与字段补齐。
//
// 为什么 HTTP 节点要单独走 JSON map 这条路（而不是和其它节点一样在 cleanNode 里用 *vo.Node 改）：
//  1. 前端初始化鉴权表单时会读取 `param.type`，但 vo.Param 结构体里只有
//     name / input / left / right / variables，根本没有 type 字段。即使在 Go 侧
//     给 *vo.Node 设上，marshal 回 JSON 时也会被丢掉，前端依旧拿到 undefined。
//  2. wanwu 导出 JSON 里 auth 的实际形状（authData.bearerTokenData / basicAuthData /
//     customData.data，每个元素再带 input.value.content 等）和 vo.Auth 不完全对齐，
//     先反序列化到结构体会让多余字段被丢弃，再写回时无法保留。
//
// 因此实现选择：
//   - 调用方先 Marshal(schema) → Unmarshal 到 map[string]any，本函数直接在 map 层修改，
//     既不污染共享 VO 结构体，也能保留 wanwu 原始字段。
//   - 顶层从 schemaMap["nodes"] 进入，再递归下钻 node["blocks"]
//     （vo.Node.Blocks 的 JSON tag），保证嵌套在循环/批处理子画布里的 HTTP 节点也会被处理。
func cleanWANWUHTTPAuthInSchemaJSON(schemaMap map[string]any) {
	if schemaMap == nil {
		return
	}

	nodes, ok := schemaMap["nodes"].([]any)
	if !ok {
		return
	}
	cleanWANWUHTTPAuthInNodesSlice(nodes)
}

// cleanWANWUHTTPAuthInNodesSlice 递归处理节点列表。
//
// 顶层画布的 schema.nodes 与复合节点（循环 type="21"、批处理等）的 node.blocks
// 在 JSON 上是同构的——都是 []*vo.Node 序列化出的对象数组——所以共用同一套逻辑：
//   - 当前节点本身若是 HTTP（type="45"），交给 cleanWANWUHTTPAuthInSingleNodeMap 处理；
//   - 不论当前节点是不是复合节点，只要带有 blocks 子数组就继续下钻。
//
// 这里不限定父节点 type，是为了让任意带 blocks 的复合节点（包括将来新增的类型）
// 内部嵌套的 HTTP 节点都能被覆盖到，避免漏改。
func cleanWANWUHTTPAuthInNodesSlice(nodes []any) {
	for _, nodeAny := range nodes {
		node, ok := nodeAny.(map[string]any)
		if !ok {
			continue
		}
		cleanWANWUHTTPAuthInSingleNodeMap(node)
		// 非复合节点序列化时因 `omitempty` 不会出现 "blocks" 键，断言会直接失败跳过；
		// 复合节点（如循环）则在此处下钻进入子画布。
		if blocks, ok := node["blocks"].([]any); ok && len(blocks) > 0 {
			cleanWANWUHTTPAuthInNodesSlice(blocks)
		}
	}
}

// cleanWANWUHTTPAuthInSingleNodeMap 对单个 HTTP 节点（type="45"）做鉴权字段清理与补齐。
//
// 处理路径（仅当节点存在且 authOpen=true 时才生效）：
//
//	node.data.inputs.auth.authData.{bearerTokenData[] | basicAuthData[] | customData.data[]}
//
// 对每个 param 元素：
//   - 把 input.value.content 置空（隐藏 token / 密钥，避免随 schema 一并导出）；
//   - 把 input.type 同步到顶层 param.type（前端表单初始化依赖该字段，但 vo.Param
//     上没有此字段，所以必须在 JSON map 层手工补齐）。
//
// 沿途每一层都用类型断言软失败，遇到结构缺失就直接 return，确保对不规则导出 JSON 也不会 panic。
func cleanWANWUHTTPAuthInSingleNodeMap(node map[string]any) {
	nodeType, _ := node["type"].(string)
	if nodeType != "45" { // HTTP 请求
		return
	}

	data, ok := node["data"].(map[string]any)
	if !ok {
		return
	}
	inputs, ok := data["inputs"].(map[string]any)
	if !ok {
		return
	}
	auth, ok := inputs["auth"].(map[string]any)
	if !ok {
		return
	}

	// 没开鉴权就没有 token 需要清，直接放行；同时也跳过 param.type 补齐，
	// 因为前端只在 authOpen 分支会读这些字段。
	authOpen, _ := auth["authOpen"].(bool)
	if !authOpen {
		return
	}

	authData, ok := auth["authData"].(map[string]any)
	if !ok {
		return
	}

	// clearTokenValue 处理单个 param：清空敏感值 + 把 input.type 抬到顶层 param.type。
	clearTokenValue := func(param map[string]any) {
		if param == nil {
			return
		}
		input, ok := param["input"].(map[string]any)
		if !ok {
			return
		}

		// 前端初始化需要顶层 param.type；wanwu 导出源只在 input.type 里提供。
		// 这一步是这条链路必须走 map 而不是 *vo.Node 的根本原因。
		if inputType, ok := input["type"]; ok {
			param["type"] = inputType
		}

		value, ok := input["value"].(map[string]any)
		if !ok {
			return
		}
		// 注意：保留 value 其它字段（如 type/rawMeta），仅清掉真正的密文。
		value["content"] = ""
	}

	// processParamList 遍历 authData 下某个数组字段（如 bearerTokenData）里的所有 param。
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

	// 前端会遍历以下三类鉴权数组（只要导出 JSON 里存在相应字段）：
	//   - bearerTokenData：Bearer Token 鉴权
	//   - basicAuthData：HTTP Basic 鉴权
	//   - customData.data：自定义 Header / Query 鉴权
	processParamList(authData, "bearerTokenData")
	processParamList(authData, "basicAuthData")
	if customData, ok := authData["customData"].(map[string]any); ok {
		processParamList(customData, "data")
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
