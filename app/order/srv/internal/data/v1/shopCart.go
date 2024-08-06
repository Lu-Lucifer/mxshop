package v1

import (
	"context"
	"database/sql"
	"mxshop/app/order/srv/internal/domain/do"
	v1 "mxshop/pkg/common/meta/v1"
)

type ShopCartStore interface {
	//查询用户购物车列表
	List(ctx context.Context, userID uint64, checked bool, meta v1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error)
	//新建某一条购物车
	Create(ctx context.Context, cartItem *do.ShoppingCartDO) error
	//获取某一条购物车详情
	Get(ctx context.Context, userID, goodsID uint64) (*do.ShoppingCartDO, error)
	//更新某一条购物车数量
	UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error
	//删除某一条购物车信息
	Delete(ctx context.Context, ID uint64) error
	//清空选中状态
	ClearCheck(ctx context.Context, userID uint64) error

	//新建订单后删除购物车对应记录
	//DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []int32) error
	DeleteByGoodsIDs(ctx context.Context, txn *sql.Tx, userID uint64, goodsIDs []int32) error
}
