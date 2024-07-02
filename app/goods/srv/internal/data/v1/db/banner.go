package db

import (
	"context"
	"gorm.io/gorm"
	v1 "mxshop/app/goods/srv/internal/data/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	metav1 "mxshop/pkg/common/meta/v1"
)

type banner struct {
	db *gorm.DB
}

// 工厂模式
func newBanner(factory *mysqlFactory) *banner {
	return &banner{
		db: factory.db,
	}
}

func (b *banner) List(ctx context.Context, opts metav1.ListMeta, orderby []string) (*do.BannerList, error) {
	//TODO implement me
	panic("implement me")
}

func (b *banner) Create(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	//TODO implement me
	panic("implement me")
}

func (b *banner) Update(ctx context.Context, txn *gorm.DB, banner *do.BannerDO) error {
	//TODO implement me
	panic("implement me")
}

func (b *banner) Delete(ctx context.Context, ID uint64) error {
	//TODO implement me
	panic("implement me")
}

var _ v1.BannerStore = &banner{}
