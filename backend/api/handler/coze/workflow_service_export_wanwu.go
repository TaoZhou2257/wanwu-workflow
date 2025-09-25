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
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	// create
	createReq := workflow.CreateWorkflowRequest{
		SpaceID: req.SpaceID,
		Name:    req.Name,
		Desc:    req.Desc,
		IconURI: "default_icon/default_workflow_icon.png",
	}
	resp, err := appworkflow.SVC.CreateWorkflow(ctx, &createReq)
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
}

// ExportWorkFlow .
// @router /api/workflow_api/export [POST]
func ExportWorkFlow(ctx context.Context, c *app.RequestContext) {
	var err error
	var req workflow.GetCanvasInfoRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		invalidParamRequestResponse(c, err.Error())
		return
	}
	resp, err := appworkflow.SVC.GetCanvasInfo(ctx, &req)
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
	if err := cleanWorkflowSchema(&schema); err != nil {
		internalServerErrorResponse(ctx, c, fmt.Errorf("failed to clean workflow schema: %v", err))
		return
	}
	var schemaStr string
	if schemaStr, err = sonic.MarshalString(schema); err != nil {
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
	case "1006": // 知识库检索节点
		cleanKnowledgeNode(node)
	case "1009": // MCP节点
		cleanMCPNode(node)
	case "1010": // GUI智能体节点
		cleanGUINode(node)
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
