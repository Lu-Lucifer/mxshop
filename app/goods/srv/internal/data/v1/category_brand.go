package v1

import (
	"context"
	"gorm.io/gorm"
	"mxshop/app/goods/srv/internal/domain/do"
	metav1 "mxshop/pkg/common/meta/v1"
)

type GoodsCategoryBrandStore interface {
	List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error)
	Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error
	Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error
	Delete(ctx context.Context, ID uint64) error
}
