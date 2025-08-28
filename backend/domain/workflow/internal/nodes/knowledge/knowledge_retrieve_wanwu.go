/*
 * author wangliang
 */

package knowledge

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/coze-dev/coze-studio/backend/application/base/ctxutil"
	http_client "github.com/coze-dev/coze-studio/backend/pkg/http-client"
	"github.com/coze-dev/coze-studio/backend/pkg/logs"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"

	"github.com/spf13/cast"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/canvas/convert"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
)

const (
	successCode  = 0
	wanWuOutput  = "output"
	ragSearchUrl = "http://rag-wanwu:8681/rag/search-knowledge-base" // todo 变成环境变量
)

type WanWuRetrieveConfig struct {
	KnowledgeIDs   []string
	RetrieveParams *RetrieveParams
}

type RetrieveParams struct {
	MatchType                   string  `json:"matchType"  validate:"required"` //matchType：vector（向量检索）、text（文本检索）、mix（混合检索：向量+文本）
	RerankModelId               string  `json:"rerankModelId"`                  //rerank模型id
	PriorityMatch               int     `json:"priorityMatch"`                  // 权重匹配，只有在混合检索模式下，选择权重设置后，这个才设置为1
	SemanticsPriority           float64 `json:"semanticsPriority"`              // 语义权重
	KeywordPriority             float64 `json:"keywordPriority"`                // 关键词权重
	RerankKeywordPriority       float64 `json:"rerankKeywordPriority"`          // 混合搜索的关键词权重
	RerankKeywordPrioritySwitch bool    `json:"rerankKeywordPrioritySwitch"`    // 混合搜索的关键词权重开关
	TopK                        int64   `json:"topK"`                           //topK 获取最高的几行
	Threshold                   float64 `json:"threshold"`                      //threshold 过滤分数阈值
	Rewrite                     bool    `json:"rewrite"`                        //是否开启重写
}

type HitParams struct {
	UserId                string        `json:"userId"`
	Question              string        `json:"question" validate:"required"`
	KnowledgeBase         []string      `json:"knowledgeBase" validate:"required"`
	Threshold             float64       `json:"threshold"`
	TopK                  int64         `json:"topK"`
	RerankModelId         string        `json:"rerank_model_id"`         // rerankId
	RerankMod             string        `json:"rerank_mod"`              // rerank_model:重排序模式，weighted_score：权重搜索
	RetrieveMethod        string        `json:"retrieve_method"`         // hybrid_search:混合搜索， semantic_search:向量搜索， full_text_search：文本搜索
	Weight                *WeightParams `json:"weights"`                 // 权重搜索下的权重配置
	RewriteQuery          bool          `json:"rewrite_query"`           // 查询重写
	ReturnMeta            bool          `json:"return_meta"`             // 展示角标
	TermWeightCoefficient *float64      `json:"term_weight_coefficient"` // 展示角标
}

type WeightParams struct {
	VectorWeight float64 `json:"vector_weight"` //语义权重
	TextWeight   float64 `json:"text_weight"`   //关键字权重
}

func (r *WanWuRetrieveConfig) Adapt(_ context.Context, n *vo.Node, _ ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	ns := &schema.NodeSchema{
		Key:     vo.NodeKey(n.ID),
		Type:    entity.NodeTypeWanWuKnowledgeRetriever,
		Name:    n.Data.Meta.Title,
		Configs: r,
	}

	inputs := n.Data.Inputs
	datasetListInfoParam := inputs.DatasetParam[0]
	datasetIDs := datasetListInfoParam.Input.Value.Content.([]any)
	knowledgeIDs := make([]string, 0, len(datasetIDs))
	for _, id := range datasetIDs {
		k := cast.ToString(id)
		knowledgeIDs = append(knowledgeIDs, k)
	}
	r.KnowledgeIDs = knowledgeIDs

	retrieveParams := &RetrieveParams{}

	var getDesignatedParamContent = func(name string) (any, bool) {
		for _, param := range inputs.DatasetParam {
			if param.Name == name {
				return param.Input.Value.Content, true
			}
		}
		return nil, false
	}

	if content, ok := getDesignatedParamContent("topK"); ok {
		topK, err := cast.ToInt64E(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.TopK = topK
	}

	if content, ok := getDesignatedParamContent("threshold"); ok {
		threshold, err := cast.ToFloat64E(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.Threshold = threshold
	}

	//vector（向量检索）、text（文本检索）、mix_rerank(混合 rerank模型)、mix_priority(混合权重)
	if content, ok := getDesignatedParamContent("matchType"); ok {
		matchType := cast.ToString(content)
		retrieveParams.MatchType = matchType
		if matchType == "mix_priority" {
			retrieveParams.PriorityMatch = 1
		}
	}

	if content, ok := getDesignatedParamContent("rerankModelId"); ok {
		retrieveParams.RerankModelId = cast.ToString(content)
	}

	if content, ok := getDesignatedParamContent("semanticsPriority"); ok {
		semanticsPriority, err := cast.ToFloat64E(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.SemanticsPriority = semanticsPriority
	}

	if content, ok := getDesignatedParamContent("keywordPriority"); ok {
		keywordPriority, err := cast.ToFloat64E(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.KeywordPriority = keywordPriority
	}

	if content, ok := getDesignatedParamContent("rerankKeywordPriority"); ok {
		rerankKeywordPriority, err := cast.ToFloat64E(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.RerankKeywordPriority = rerankKeywordPriority
	}

	if content, ok := getDesignatedParamContent("rerankKeywordPrioritySwitch"); ok {
		rerankKeywordPrioritySwitch, err := cast.ToBoolE(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.RerankKeywordPrioritySwitch = rerankKeywordPrioritySwitch
	}

	if content, ok := getDesignatedParamContent("rewrite"); ok {
		rewrite, err := cast.ToBoolE(content)
		if err != nil {
			return nil, err
		}
		retrieveParams.Rewrite = rewrite
	}

	r.RetrieveParams = retrieveParams

	if err := convert.SetInputsForNodeSchema(n, ns); err != nil {
		return nil, err
	}

	if err := convert.SetOutputTypesForNodeSchema(n, ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (r *WanWuRetrieveConfig) Build(_ context.Context, _ *schema.NodeSchema, _ ...schema.BuildOption) (any, error) {
	if len(r.KnowledgeIDs) == 0 {
		return nil, errors.New("knowledge ids are required")
	}

	if r.RetrieveParams == nil {
		return nil, errors.New("retrieval params is required")
	}

	return &WanWuRetrieve{
		knowledgeIDs:   r.KnowledgeIDs,
		retrieveParams: r.RetrieveParams,
	}, nil
}

type WanWuRetrieve struct {
	knowledgeIDs   []string
	retrieveParams *RetrieveParams
}

type RagKnowledgeHitResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    *HitData `json:"data"`
}

type HitData struct {
	Prompt     string `json:"prompt"`
	SearchList []any  `json:"searchList"`
	Score      []any  `json:"score"` //实际是[]float
}

// ChunkSearchList HitData 的SearchList
type ChunkSearchList struct {
	Title    string      `json:"title"`
	Snippet  string      `json:"snippet"`
	KbName   string      `json:"kb_name"`
	MetaData interface{} `json:"meta_data"`
}

func (kr *WanWuRetrieve) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	query, ok := input["Query"].(string)
	if !ok {
		return nil, errors.New("capital query key is required")
	}

	userId := ctxutil.MustGetUIDFromCtx(ctx)
	userIdStr := strconv.Itoa(int(userId))

	retrieveParams := kr.retrieveParams
	priorityMatch := retrieveParams.PriorityMatch

	var termWeightCoefficient *float64 = nil
	if retrieveParams.RerankKeywordPrioritySwitch {
		termWeightCoefficient = &retrieveParams.RerankKeywordPriority
	}
	req := &HitParams{
		Question:              query,
		KnowledgeBase:         kr.knowledgeIDs,
		TopK:                  retrieveParams.TopK,
		Threshold:             retrieveParams.Threshold,
		RerankModelId:         buildRerankId(priorityMatch, retrieveParams.RerankModelId),
		RetrieveMethod:        buildRetrieveMethod(retrieveParams.MatchType),
		RerankMod:             buildRerankMod(priorityMatch),
		Weight:                buildWeight(priorityMatch, retrieveParams.SemanticsPriority, retrieveParams.KeywordPriority),
		RewriteQuery:          retrieveParams.Rewrite,
		UserId:                userIdStr,
		ReturnMeta:            true,
		TermWeightCoefficient: termWeightCoefficient,
	}

	response, err := ragKnowledgeSearch(ctx, req)
	if err != nil {
		return nil, err
	}
	result := make(map[string]any)
	result[wanWuOutput] = map[string]any{
		"prompt":     response.Data.Prompt,
		"score":      response.Data.Score,
		"searchList": response.Data.SearchList,
	}

	return result, nil
}

// ragKnowledgeSearch rag命中测试
func ragKnowledgeSearch(ctx context.Context, knowledgeHitParams *HitParams) (*RagKnowledgeHitResp, error) {
	paramsByte, err := sonic.Marshal(knowledgeHitParams)
	if err != nil {
		logs.CtxErrorf(ctx, "ragKnowledgeSearch params marsh error %v", err)
		return nil, err
	}
	result, err := http_client.GetDefaultClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        ragSearchUrl,
		Body:       paramsByte,
		Timeout:    time.Duration(10) * time.Second,
		MonitorKey: "rag_knowledge_hit",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp RagKnowledgeHitResp
	if err := sonic.Unmarshal(result, &resp); err != nil {
		logs.CtxErrorf(ctx, "ragKnowledgeSearch result Unmarshal error %v", err)
		return nil, err
	}
	if resp.Code != successCode {
		return nil, errors.New(resp.Message)
	}
	return &resp, nil
}

// buildRerankId 构造重排序模型id
func buildRerankId(priorityType int, rerankId string) string {
	if priorityType == 1 {
		return ""
	}
	return rerankId
}

// buildRetrieveMethod 构造检索方式
func buildRetrieveMethod(matchType string) string {
	switch matchType {
	case "vector":
		return "semantic_search"
	case "text":
		return "full_text_search"
	case "mix_priority":
	case "mix_rerank":
		return "hybrid_search"
	}
	return ""
}

// buildRerankMod 构造重排序模式
func buildRerankMod(priorityType int) string {
	if priorityType == 1 {
		return "weighted_score"
	}
	return "rerank_model"
}

// buildWeight 构造权重信息
func buildWeight(priorityType int, semanticsPriority float64, keywordPriority float64) *WeightParams {
	if priorityType != 1 {
		return nil
	}
	return &WeightParams{
		VectorWeight: semanticsPriority,
		TextWeight:   keywordPriority,
	}
}
