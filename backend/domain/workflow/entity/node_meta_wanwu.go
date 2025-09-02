package entity

// Wanwu Node Type Definition
const (
	NodeTypeWanwuIntentDetector     NodeType = "WanwuIntentDetector"
	NodeTypeWanWuKnowledgeRetriever NodeType = "WanWuKnowledgeRetriever"
	NodeTypeWanWuFileGenerator      NodeType = "WanWuFileGenerator"
	NodeTypeWanWuFileParser         NodeType = "WanWuFileParser"
	NodeTypeWanWuMCPTool            NodeType = "WanWuMCPTool"
)

// Wanwu NodeTypeMetas Init
func init() {
	// wanwu新增节点分类
	Categories = append(Categories, Category{
		Key:      "document",
		Name:     "文档",
		EnUSName: "Document",
	})

	// wanwu禁用一些节点
	NodeTypeMetas[NodeTypePlugin].Disabled = true
	NodeTypeMetas[NodeTypeKnowledgeRetriever].Disabled = true
	NodeTypeMetas[NodeTypeSubWorkflow].Disabled = true
	NodeTypeMetas[NodeTypeDatabaseCustomSQL].Disabled = true
	NodeTypeMetas[NodeTypeQuestionAnswer].Disabled = true
	NodeTypeMetas[NodeTypeKnowledgeIndexer].Disabled = true
	NodeTypeMetas[NodeTypeMessageList].Disabled = true
	NodeTypeMetas[NodeTypeClearMessage].Disabled = true
	NodeTypeMetas[NodeTypeCreateConversation].Disabled = true
	NodeTypeMetas[NodeTypeVariableAssigner].Disabled = true
	NodeTypeMetas[NodeTypeDatabaseUpdate].Disabled = true
	NodeTypeMetas[NodeTypeDatabaseQuery].Disabled = true
	NodeTypeMetas[NodeTypeDatabaseDelete].Disabled = true
	NodeTypeMetas[NodeTypeDatabaseInsert].Disabled = true
	NodeTypeMetas[NodeTypeKnowledgeDeleter].Disabled = true
	// 和前端约定，反序列化节点ID 59 -> 1059
	NodeTypeMetas[NodeTypeJsonDeserialization].ID = 1059

	// wanwu新增节点
	NodeTypeMetas[NodeTypeWanWuKnowledgeRetriever] = &NodeTypeMeta{
		ID:           1006,
		Key:          NodeTypeWanWuKnowledgeRetriever,
		DisplayKey:   "Dataset",
		Name:         "知识库检索",
		Category:     "data",
		Desc:         "在选定的知识中,根据输入变量召回最匹配的信息,并以列表形式返回",
		Color:        "#FF811A",
		IconURL:      "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeQuery-v2.jpg",
		SupportBatch: false,
		ExecutableMeta: ExecutableMeta{
			PreFillZero: true,
			PostFillNil: true,
		},
		EnUSName:        "Knowledge retrieval",
		EnUSDescription: "In the selected knowledge, the best matching information is recalled based on the input variable and returned as an Array.",
	}

	NodeTypeMetas[NodeTypeWanWuFileGenerator] = &NodeTypeMeta{
		ID:           1007,
		Key:          NodeTypeWanWuFileGenerator,
		DisplayKey:   "FileGenerator",
		Name:         "文档生成",
		Category:     "document",
		Desc:         "输入文档的内容、格式和文件名，可以生成文档下载链接",
		Color:        "#FF811A",
		IconURL:      "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeQuery-v2.jpg",
		SupportBatch: false,
		ExecutableMeta: ExecutableMeta{
			PreFillZero: true,
			PostFillNil: true,
		},
		EnUSName:        "File generator",
		EnUSDescription: "Generate documents using both the document title and content.",
	}

	NodeTypeMetas[NodeTypeWanWuFileParser] = &NodeTypeMeta{
		ID:           1008,
		Key:          NodeTypeWanWuFileParser,
		DisplayKey:   "FileParser",
		Name:         "文档解析",
		Category:     "document",
		Desc:         "输入txt、pdf、docx、xlsx、csv、pptx等格式文档的URL，可以解析提取出文档的文本内容",
		Color:        "#FF811A",
		IconURL:      "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeQuery-v2.jpg",
		SupportBatch: false,
		ExecutableMeta: ExecutableMeta{
			PreFillZero: true,
			PostFillNil: true,
		},
		EnUSName:        "File generator",
		EnUSDescription: "Parse documents content.",
	}

	NodeTypeMetas[NodeTypeWanWuMCPTool] = &NodeTypeMeta{
		ID:           1009,
		Key:          NodeTypeWanWuMCPTool,
		DisplayKey:   "MCPTool",
		Name:         "MCP工具",
		Category:     "utilities",
		Desc:         "用于调用MCP服务的工具",
		Color:        "#FF811A",
		IconURL:      "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-KnowledgeQuery-v2.jpg",
		SupportBatch: false,
		ExecutableMeta: ExecutableMeta{
			PreFillZero: true,
			PostFillNil: true,
		},
		EnUSName:        "MCP tool",
		EnUSDescription: "Used to call MCP tools.",
	}

	// 只用于示例，并不开放（Disable = true）
	NodeTypeMetas[NodeTypeWanwuIntentDetector] = &NodeTypeMeta{
		Disabled:     true,
		ID:           1022,
		Key:          NodeTypeWanwuIntentDetector,
		DisplayKey:   "Intent",
		Name:         "意图识别",
		Category:     "logic",
		Desc:         "用于用户输入的意图识别，并将其与预设意图选项进行匹配。",
		Color:        "#00B2B2",
		IconURL:      "https://lf3-static.bytednsdoc.com/obj/eden-cn/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Intent-v2.jpg",
		SupportBatch: false,
		ExecutableMeta: ExecutableMeta{
			PreFillZero:     true,
			PostFillNil:     true,
			MayUseChatModel: true,
		},
		EnUSName:        "Intent recognition",
		EnUSDescription: "Used for recognizing the intent in user input and matching it with preset intent options.",
	}
}
