package memory

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/api/model/base"
	"github.com/coze-dev/coze-studio/backend/api/model/data/database/table"
	"github.com/coze-dev/coze-studio/backend/types/consts"
)

func (d *DatabaseApplicationService) GetConnectorNameByWanwu(ctx context.Context, req *table.GetSpaceConnectorListRequest) (*table.GetSpaceConnectorListResponse, error) {
	return &table.GetSpaceConnectorListResponse{
		ConnectorList: []*table.ConnectorInfo{
			{
				ConnectorID:   consts.CozeConnectorID,
				ConnectorName: "Web",
			},
			{
				ConnectorID:   consts.APIConnectorID,
				ConnectorName: "API",
			},
		},

		Code: 0,
		Msg:  "success",
		BaseResp: &base.BaseResp{
			StatusCode:    0,
			StatusMessage: "success",
		},
	}, nil
}
