package v1

import (
	"context"
	gpb "mxshop/api/goods/v1"
	"mxshop/app/mxshop/api/internal/data"
)

type GoodsSrv interface {
	List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error)
}

type goodsService struct {
	data data.DataFactory
}

func NewGoods(data data.DataFactory) *goodsService {
	return &goodsService{
		data: data,
	}
}

func (gs *goodsService) List(ctx context.Context, request *gpb.GoodsFilterRequest) (*gpb.GoodsListResponse, error) {
	return gs.data.Goods().GoodsList(ctx, request)

}

var _ GoodsSrv = &goodsService{}
