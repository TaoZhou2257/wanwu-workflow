package search

import (
	"context"

	"github.com/coze-dev/coze-studio/backend/domain/search/entity"
)

type ResourceEventBusImpl struct{}

func (r *ResourceEventBusImpl) PublishResources(ctx context.Context, event *entity.ResourceDomainEvent) error {
	return nil
}
