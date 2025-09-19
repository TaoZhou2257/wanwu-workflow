package router

import (
	wanwu_mock "github.com/UnicomAI/wanwu-workflow/wanwu-backend/internal/server/http/handler/wanwu-mock"
	"github.com/cloudwego/hertz/pkg/app"
	hertz_server "github.com/cloudwego/hertz/pkg/app/server"
	"github.com/coze-dev/coze-studio/backend/api/handler/coze"
)

// ！！！同步于：backend/api/router/coze/api.go
// ！！！同步于：backend/api/router/coze/api.go
// ！！！同步于：backend/api/router/coze/api.go

func Register(r *hertz_server.Hertz) {
	// TODO auto generated

	// FIXME 这里实际对应的是 ./configs/static/api/static 而非 ./configs/static ??
	r.Static("/api/static", "./configs/static")

	root := r.Group("/", rootMw()...)
	{
		_api := root.Group("/api", _apiMw()...)
		{
			_bot := _api.Group("/bot", _botMw()...)
			_bot.POST("/get_type_list", append(_gettypelistMw(), coze.GetTypeList)...)
		}
		{
			_common := _api.Group("/common", _commonMw()...)
			{
				_upload := _common.Group("/upload", _uploadMw()...)
				_upload.GET("/apply_upload_action", append(_applyuploadactionMw(), coze.ApplyUploadActionByWanwu)...)
				_upload.POST("/apply_upload_action", append(_applyuploadaction0Mw(), coze.ApplyUploadActionByWanwu)...)
				_upload.POST("/*tos_uri", append(_commonuploadMw(), coze.CommonUpload)...)
			}
		}
		{
			_playground_api := _api.Group("/playground_api", _playground_apiMw()...)
			{
				_space := _playground_api.Group("/space", _spaceMw()...)
				_space.POST("/list", append(_getspacelistv2Mw(), wanwu_mock.GetSpaceListV2)...)
			}
		}
		{
			_workflow_api := _api.Group("/workflow_api", _workflow_apiMw()...)
			_workflow_api.GET("/apiDetail", append(_getapidetailMw(), coze.GetApiDetail)...)
			_workflow_api.POST("/batch_delete", append(_batchdeleteworkflowMw(), coze.BatchDeleteWorkflow)...)
			_workflow_api.POST("/cancel", append(_cancelworkflowMw(), coze.CancelWorkFlow)...)
			_workflow_api.POST("/canvas", append(_getcanvasinfoMw(), coze.GetCanvasInfo)...)
			_workflow_api.POST("/copy", append(_copyworkflowMw(), coze.CopyWorkflow)...)
			_workflow_api.POST("/copy_wk_template", append(_copywktemplateapiMw(), coze.CopyWkTemplateApi)...)
			_workflow_api.POST("/create", append(_createworkflowMw(), coze.CreateWorkflow)...)
			_workflow_api.POST("/delete", append(_deleteworkflowMw(), coze.DeleteWorkflow)...)
			_workflow_api.POST("/delete_strategy", append(_getdeletestrategyMw(), coze.GetDeleteStrategy)...)
			_workflow_api.POST("/example_workflow_list", append(_getexampleworkflowlistMw(), coze.GetExampleWorkFlowListByWanwu)...)
			_workflow_api.GET("/get_node_execute_history", append(_getnodeexecutehistoryMw(), coze.GetNodeExecuteHistory)...)
			_workflow_api.GET("/get_process", append(_getworkflowprocessMw(), coze.GetWorkFlowProcess)...)
			_workflow_api.POST("/get_trace", append(_gettracesdkMw(), coze.GetTraceSDK)...)
			_workflow_api.POST("/history_schema", append(_gethistoryschemaMw(), coze.GetHistorySchema)...)
			_workflow_api.POST("/list_publish_workflow", append(_listpublishworkflowMw(), coze.ListPublishWorkflow)...)
			_workflow_api.POST("/list_spans", append(_listrootspansMw(), coze.ListRootSpans)...)
			_workflow_api.POST("/llm_fc_setting_detail", append(_getllmnodefcsettingdetailMw(), coze.GetLLMNodeFCSettingDetail)...)
			_workflow_api.POST("/llm_fc_setting_merged", append(_getllmnodefcsettingsmergedMw(), coze.GetLLMNodeFCSettingsMerged)...)
			_workflow_api.POST("/nodeDebug", append(_workflownodedebugv2Mw(), coze.WorkflowNodeDebugV2)...)
			_workflow_api.POST("/node_panel_search", append(_nodepanelsearchMw(), coze.NodePanelSearch)...)
			_workflow_api.POST("/node_template_list", append(_nodetemplatelistMw(), wanwu_mock.NodeTemplateListByWanwu)...)
			_workflow_api.POST("/node_type", append(_queryworkflownodetypesMw(), coze.QueryWorkflowNodeTypes)...)
			_workflow_api.POST("/publish", append(_publishworkflowMw(), coze.PublishWorkflow)...)
			_workflow_api.POST("/released_workflows", append(_getreleasedworkflowsMw(), coze.GetReleasedWorkflows)...)
			_workflow_api.POST("/save", append(_saveworkflowMw(), coze.SaveWorkflow)...)
			_workflow_api.POST("/sign_image_url", append(_signimageurlMw(), coze.SignImageURL)...)
			_workflow_api.POST("/test_resume", append(_workflowtestresumeMw(), coze.WorkFlowTestResume)...)
			_workflow_api.POST("/test_run", append(_workflowtestrunMw(), coze.WorkFlowTestRun)...)
			_workflow_api.POST("/update_meta", append(_updateworkflowmetaMw(), coze.UpdateWorkflowMeta)...)
			_workflow_api.POST("/validate_tree", append(_validatetreeMw(), coze.ValidateTree)...)
			_workflow_api.POST("/workflow_detail", append(_getworkflowdetailMw(), coze.GetWorkflowDetail)...)
			_workflow_api.POST("/workflow_detail_info", append(_getworkflowdetailinfoMw(), coze.GetWorkflowDetailInfoByWanwu)...)
			_workflow_api.POST("/workflow_list", append(_getworkflowlistMw(), coze.GetWorkFlowList)...)
			_workflow_api.POST("/workflow_references", append(_getworkflowreferencesMw(), coze.GetWorkflowReferences)...)
			{
				_chat_flow_role := _workflow_api.Group("/chat_flow_role", _chat_flow_roleMw()...)
				_chat_flow_role.POST("/create", append(_createchatflowroleMw(), coze.CreateChatFlowRole)...)
				_chat_flow_role.POST("/delete", append(_deletechatflowroleMw(), coze.DeleteChatFlowRole)...)
				_chat_flow_role.GET("/get", append(_getchatflowroleMw(), coze.GetChatFlowRole)...)
			}
			{
				_project_conversation := _workflow_api.Group("/project_conversation", _project_conversationMw()...)
				_project_conversation.POST("/create", append(_createprojectconversationdefMw(), coze.CreateProjectConversationDef)...)
				_project_conversation.POST("/delete", append(_deleteprojectconversationdefMw(), coze.DeleteProjectConversationDef)...)
				_project_conversation.GET("/list", append(_listprojectconversationdefMw(), coze.ListProjectConversationDef)...)
				_project_conversation.POST("/update", append(_updateprojectconversationdefMw(), coze.UpdateProjectConversationDef)...)
			}
			{
				_upload1 := _workflow_api.Group("/upload", _upload1Mw()...)
				_upload1.POST("/auth_token", append(_getworkflowuploadauthtokenMw(), coze.GetWorkflowUploadAuthToken)...)
			}
			_plugin_api := _api.Group("/plugin_api", _plugin_apiMw()...)
			_plugin_api.POST("/get_playground_plugin_list", append(_getplaygroundpluginlistMw(), coze.GetPlaygroundPluginListByWanwu)...)

			// --- wanwu adapt ---
			_workflow_api.POST("/workflow_list_by_wanwu", []app.HandlerFunc{coze.GetWorkFlowListByWanwu}...)
			_workflow_api.POST("/workflow_select_by_wanwu", []app.HandlerFunc{coze.GetWorkFlowSelectByWanwu}...)

		}
	}
	{
		_v1 := root.Group("/v1", _v1Mw()...)
		{
			_workflow := _v1.Group("/workflow", _workflowMw()...)

			// --- wanwu adapt ---
			_workflow.POST("/list_schema_by_wanwu", []app.HandlerFunc{coze.ListWorkFlowOpenAPIV3SchemaByWanwu}...)
			_workflow.GET("/:workflow_id/schema_by_wanwu", []app.HandlerFunc{coze.GetWorkFlowOpenAPIV3SchemaByWanwu}...)
			_workflow.POST("/:workflow_id/run_by_wanwu", append(_openapirunflowMw(), coze.OpenAPIRunWorkFlowByWanwu)...)
		}
	}
}
