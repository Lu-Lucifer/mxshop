package v1

import (
	"context"
	proto "mxshop/api/goods/v1"
	"mxshop/app/goods/srv/internal/domain/do"
)

type GoodsFilterRequest struct {
	*proto.GoodsFilterRequest
	CategoryIDs []interface{}
}

type GoodsStore interface {
	Create(ctx context.Context, goods *do.GoodsSearchDO) error
	Delete(ctx context.Context, ID uint64) error
	Update(ctx context.Context, goods *do.GoodsSearchDO) error
	//es条件搜索
	Search(ctx context.Context, request *GoodsFilterRequest) (*do.GoodsSearchDOList, error)
}
