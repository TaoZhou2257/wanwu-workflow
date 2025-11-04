package impl

import (
	"context"
	"fmt"
	"path"

	"github.com/coze-dev/coze-studio/backend/api/model/app/developer_api"
	crossupload "github.com/coze-dev/coze-studio/backend/crossdomain/upload"
	"github.com/coze-dev/coze-studio/backend/infra/storage"
	"github.com/google/uuid"
)

type wanwuUploader struct {
	oss storage.Storage
}

func NewWanwuUploader(oss storage.Storage) crossupload.WanwuUploader {
	return &wanwuUploader{
		oss: oss,
	}
}

func (s *wanwuUploader) UploadFileByByte(ctx context.Context, fileName string, data []byte) (*crossupload.UploadFileResp, error) {
	fileExt := path.Ext(fileName)

	var BizType developer_api.FileBizType = 0
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	objectName := fmt.Sprintf("%s/%s", BizType.String(), newFileName)

	err := s.oss.PutObject(ctx, objectName, data)
	if err != nil {
		return nil, err
	}

	url, err := s.oss.GetObjectUrl(ctx, objectName)
	if err != nil {
		return nil, err
	}

	return &crossupload.UploadFileResp{
		URL: url,
		URI: objectName,
	}, nil
}
