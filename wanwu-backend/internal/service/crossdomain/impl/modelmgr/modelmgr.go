package modelmgr

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/infra/modelmgr"
)

var defaultMock *mock = &mock{}

func DefaultMock() *mock {
	return defaultMock
}

type mock struct{}

func (m *mock) ListModel(ctx context.Context, req *modelmgr.ListModelRequest) (*modelmgr.ListModelResponse, error) {
	return &modelmgr.ListModelResponse{}, nil
}

func (m *mock) ListInUseModel(ctx context.Context, limit int, cursor *string) (*modelmgr.ListModelResponse, error) {
	panic("ListInUseModel not implemented")
}

func (m *mock) MGetModelByID(ctx context.Context, req *modelmgr.MGetModelRequest) ([]*modelmgr.Model, error) {
	panic("MGetModelByID not implemented")
}
