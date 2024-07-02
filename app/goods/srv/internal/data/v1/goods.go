package v1

import (
	"context"
	"gorm.io/gorm"
	"mxshop/app/goods/srv/internal/domain/do"
	metav1 "mxshop/pkg/common/meta/v1"
)

type GoodsStore interface {
	//根据id查找商品信息
	Get(ctx context.Context, ID uint64) (*do.GoodsDO, error)
	ListByIDs(ctx context.Context, ids []uint64, orderby []string) (*do.GoodsDOList, error)
	//条件搜索
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*do.GoodsDOList, error)
	//第一种方案
	Create(ctx context.Context, goods *do.GoodsDO) error
	//第二种方案
	CreateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Update(ctx context.Context, goods *do.GoodsDO) error
	UpdateInTxn(ctx context.Context, txn *gorm.DB, goods *do.GoodsDO) error
	Delete(ctx context.Context, ID uint64) error
	DeleteInTxn(ctx context.Context, txn *gorm.DB, ID uint64) error
	//开启事务
	Begin() *gorm.DB
}
