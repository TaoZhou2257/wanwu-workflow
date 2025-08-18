package search

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/search/entity"
)

var defaultResourceEventBusMock *resourceEventBusMock = &resourceEventBusMock{}

func DefaultResourceEventBusMock() *resourceEventBusMock {
	return defaultResourceEventBusMock
}

type resourceEventBusMock struct{}

func (r *resourceEventBusMock) PublishResources(ctx context.Context, event *entity.ResourceDomainEvent) error {
	return nil
}
