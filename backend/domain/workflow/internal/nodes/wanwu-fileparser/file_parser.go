/*
 * author wangliang
 */

package wanwu_fileparser

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/canvas/convert"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
	http_client "github.com/coze-dev/coze-studio/backend/pkg/http-client"
	"github.com/coze-dev/coze-studio/backend/pkg/sonic"
)

type WanWuRetrieveConfig struct {
}

func (r *WanWuRetrieveConfig) Adapt(_ context.Context, n *vo.Node, _ ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	ns := &schema.NodeSchema{
		Key:     vo.NodeKey(n.ID),
		Type:    entity.NodeTypeWanWuFileParser,
		Name:    n.Data.Meta.Title,
		Configs: r,
	}

	if err := convert.SetInputsForNodeSchema(n, ns); err != nil {
		return nil, err
	}

	if err := convert.SetOutputTypesForNodeSchema(n, ns); err != nil {
		return nil, err
	}

	return ns, nil
}

func (r *WanWuRetrieveConfig) Build(_ context.Context, _ *schema.NodeSchema, _ ...schema.BuildOption) (any, error) {
	return &WanWuRetrieve{}, nil
}

type WanWuRetrieve struct {
}

type FileParserParams struct {
	UploadFileUrl string `json:"upload_file_url"`
}

type FileParserResp struct {
	Metadata *FileParserMeta `json:"metadata"`
	Text     string          `json:"text"`
}

type FileParserMeta struct {
	FileName string `json:"file_name"`
}

func (kr *WanWuRetrieve) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	fileUrl, ok := input["FileUrl"].(string)
	if !ok {
		return nil, errors.New("capital query key is required")
	}

	req := &FileParserParams{
		UploadFileUrl: fileUrl,
	}

	response, err := fileParser(ctx, req)
	if err != nil {
		return nil, err
	}

	//循环拼接结果
	var resultText strings.Builder
	for _, resp := range response {
		resultText.WriteString(resp.Text)
	}
	result := map[string]any{
		"text": resultText.String(),
	}

	return result, nil
}

// fileParser 文档解析
func fileParser(ctx context.Context, fileParserParams *FileParserParams) ([]*FileParserResp, error) {
	paramsByte, err := sonic.Marshal(fileParserParams)
	if err != nil {
		return nil, err
	}
	result, err := http_client.GetDefaultClient().PostJson(ctx, &http_client.HttpRequestParams{
		Url:        os.Getenv("WANWU_FILE_PARSER_URL"),
		Body:       paramsByte,
		Timeout:    time.Duration(10) * time.Second,
		MonitorKey: "file_parser",
		LogLevel:   http_client.LogAll,
	})
	if err != nil {
		return nil, err
	}
	var resp []*FileParserResp
	if err := sonic.Unmarshal(result, &resp); err != nil {
		//log.Errorf(err.Error())
		return nil, err
	}
	if len(resp) == 0 {
		return nil, errors.New("file_parser error")
	}
	return resp, nil
}
