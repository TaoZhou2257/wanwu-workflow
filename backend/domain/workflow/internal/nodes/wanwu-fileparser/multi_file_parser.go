/*
 * author wangliang
 */

package wanwu_fileparser

import (
	"context"
	"errors"
	"fmt"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/entity/vo"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/canvas/convert"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/nodes"
	"github.com/coze-dev/coze-studio/backend/domain/workflow/internal/schema"
	"github.com/spf13/cast"
	"sync"
)

const (
	multiFileCount = 10
)

type WanWuMultiFileParserConfig struct {
}

func (r *WanWuMultiFileParserConfig) Adapt(_ context.Context, n *vo.Node, _ ...nodes.AdaptOption) (*schema.NodeSchema, error) {
	ns := &schema.NodeSchema{
		Key:     vo.NodeKey(n.ID),
		Type:    entity.NodeTypeWanWuMultiFileParser,
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

func (r *WanWuMultiFileParserConfig) Build(_ context.Context, _ *schema.NodeSchema, _ ...schema.BuildOption) (any, error) {
	return &WanWuMultiFileParser{}, nil
}

type WanWuMultiFileParser struct {
}

type FileParseResult struct {
	Text  string `json:"text"`
	Error error  `json:"error"`
}

func (kr *WanWuMultiFileParser) Invoke(ctx context.Context, input map[string]any) (map[string]any, error) {
	fileUrlList, ok := input["FileUrlList"].([]any)
	if !ok {
		return nil, errors.New("capital FileUrlList key is required")
	}
	var fileUrls []string
	for _, fileUrl := range fileUrlList {
		k := cast.ToString(fileUrl)
		fileUrls = append(fileUrls, k)
	}
	if len(fileUrls) > multiFileCount {
		return nil, errors.New("multi file parse most ten file")
	}
	parserResult, err := batchFileParser(ctx, fileUrls)
	if err != nil {
		return nil, err
	}
	result := map[string]any{
		"text": parserResult,
	}
	return result, nil
}

func batchFileParser(ctx context.Context, fileList []string) ([]any, error) {
	var wg = sync.WaitGroup{}
	var parserResultList []*FileParseResult
	for _, fileUrl := range fileList {
		task := submitTask(ctx, &wg, fileUrl)
		parserResultList = append(parserResultList, task)
	}
	wg.Wait()
	var parserResult []any
	var hasSuccess = false
	for _, result := range parserResultList {
		err := result.Error
		if err == nil {
			hasSuccess = true
			parserResult = append(parserResult, result.Text)
		} else {
			parserResult = append(parserResult, err.Error())
		}
	}
	//如果没有一个文件成功，则返回错误
	if !hasSuccess {
		return nil, errors.New("file parser error")
	}
	return parserResult, nil
}

func submitTask(ctx context.Context, wg *sync.WaitGroup, fileUrl string) *FileParseResult {
	fileParseResult := &FileParseResult{}
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fileParseResult.Error = errors.New(fmt.Sprintf("panic: %v", r))
			}
			wg.Done()
		}()
		parserResult, err := CommonFileParser(ctx, fileUrl)
		fileParseResult.Text = parserResult
		fileParseResult.Error = err
	}()
	return fileParseResult
}
