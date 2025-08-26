package entity

// Wanwu Node Type Definition
const (
	NodeTypeWanwuIntentDetector NodeType = "WanwuIntentDetector"
)

// Wanwu NodeTypeMetas Init
func init() {
	NodeTypeMetas[NodeTypeWanwuIntentDetector] = &NodeTypeMeta{
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
