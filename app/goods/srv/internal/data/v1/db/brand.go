package db

import (
	"context"
	"gorm.io/gorm"
	v1 "mxshop/app/goods/srv/internal/data/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	metav1 "mxshop/pkg/common/meta/v1"
)

type brands struct {
	db *gorm.DB
}

func NewBrands(db *gorm.DB) *brands {
	return &brands{db: db}
}

func newBrands(factory *mysqlFactory) *brands {
	return &brands{
		db: factory.db,
	}
}

func (b *brands) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BrandsDOList, error) {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Create(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Update(ctx context.Context, txn *gorm.DB, brands *do.BrandsDO) error {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

func (b *brands) Get(ctx context.Context, ID uint64) (*do.BrandsDO, error) {
	//TODO implement me
	panic("implement me")
}

var _ v1.BrandsStore = &brands{}
