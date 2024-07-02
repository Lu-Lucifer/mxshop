package db

import (
	"context"
	"gorm.io/gorm"
	v1 "mxshop/app/goods/srv/internal/data/v1"
	"mxshop/app/goods/srv/internal/domain/do"
	"mxshop/app/pkg/code"
	code2 "mxshop/gmicro/code"
	"mxshop/pkg/errors"
)

type categorys struct {
	db *gorm.DB
}

func NewCategory(db *gorm.DB) *categorys {
	return &categorys{db: db}
}

func newCategory(factory *mysqlFactory) *categorys {
	return &categorys{
		db: factory.db,
	}
}

func (c *categorys) Get(ctx context.Context, ID uint64) (*do.CategoryDO, error) {
	category := &do.CategoryDO{}
	err := c.db.Preload("SubCategory").Preload("SubCategory.SubCategory").First(category, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrCategoryNotFound, "category not found")
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return category, nil
}

func (c *categorys) ListAll(ctx context.Context, orderby []string) (*do.CategoryDOList, error) {
	ret := &do.CategoryDOList{}

	//排序
	query := c.db
	for _, value := range orderby {
		query = query.Order(value)
	}
	d := query.Where("level = 1").Preload("SubCategory.SubCategory").Find(&ret.Items)
	return ret, d.Error
}

func (c *categorys) Create(ctx context.Context, category *do.CategoryDO) error {
	tx := c.db.Create(category)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *categorys) Update(ctx context.Context, category *do.CategoryDO) error {
	tx := c.db.Save(category)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (c *categorys) Delete(ctx context.Context, ID uint64) error {
	return c.db.Where("id = ?", ID).Delete(&do.CategoryDO{}).Error
}

var _ v1.CategoryStore = &categorys{}
