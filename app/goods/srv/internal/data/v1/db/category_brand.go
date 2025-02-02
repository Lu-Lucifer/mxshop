package db

import (
	"context"
	"gorm.io/gorm"
	v1 "mxshop/app/goods/srv/internal/data/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	metav1 "mxshop/pkg/common/meta/v1"
)

type categoryBrands struct {
	db *gorm.DB
}

// 工厂模式
func newCategoryBrands(factory *mysqlFactory) *categoryBrands {
	return &categoryBrands{
		db: factory.db,
	}
}

func (cb *categoryBrands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.GoodsCategoryBrandList, error) {
	//TODO implement me
	panic("implement me")
}

func (cb *categoryBrands) Create(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	//TODO implement me
	panic("implement me")
}

func (cb *categoryBrands) Update(ctx context.Context, txn *gorm.DB, gcb *do.GoodsCategoryBrandDO) error {
	//TODO implement me
	panic("implement me")
}

func (cb *categoryBrands) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

var _ v1.GoodsCategoryBrandStore = &categoryBrands{}
