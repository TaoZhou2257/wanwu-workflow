package user

import (
	"context"

	gin_util "github.com/UnicomAI/wanwu/pkg/gin-util"
	"github.com/UnicomAI/wanwu/pkg/util"
	"github.com/coze-dev/coze-studio/backend/domain/user/entity"
	"github.com/coze-dev/coze-studio/backend/pkg/ctxcache"
)

type Impl struct {
}

func (u *Impl) GetUserSpaceList(ctx context.Context, userID int64) (spaces []*entity.Space, err error) {
	orgID, ok := ctxcache.Get[string](ctx, gin_util.X_ORG_ID)
	if !ok {
		return nil, nil
	}
	return []*entity.Space{
		{ID: util.MustI64(orgID)},
	}, nil
}
