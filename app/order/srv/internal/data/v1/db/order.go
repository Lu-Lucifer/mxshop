package db

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	v1 "mxshop/app/order/srv/internal/data/v1"
	"mxshop/app/order/srv/internal/domain/do"
	code2 "mxshop/gmicro/code"
	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/errors"
)

type orders struct {
	db *gorm.DB
}

func newOrders(factory *dataFactory) *orders {
	return &orders{db: factory.db}
}
func (o *orders) Get(ctx context.Context, orderSn string) (*do.OrderInfoDO, error) {
	var order do.OrderInfoDO
	err := o.db.Preload("OrderGoods").Where("order_sn = ?", orderSn).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (o *orders) List(ctx context.Context, userID uint64, meta metav1.ListMeta, orderby []string) (*do.OrderInfoDOList, error) {
	ret := &do.OrderInfoDOList{}
	var limit, offset int
	if meta.PageSize == 0 {
		limit = 10
	} else {
		limit = meta.PageSize
	}
	if meta.Page > 0 {
		offset = (meta.Page - 1) * limit
	}
	//排序
	query := o.db.Preload("OrderGoods")
	for _, value := range orderby {
		//坑：如果db改掉了？
		//u.db=u.db.Order(value)
		query = query.Order(value)
	}
	//查询 - 发起多个请求
	d := query.Offset(offset).Limit(limit).Find(&ret.Items).Count(&ret.TotalCount)
	if d.Error != nil {
		return nil, errors.WithCode(code2.ErrDatabase, d.Error.Error())
	}
	return ret, nil

}

// 创建订单之后删除对应的购物车记录
//
//	func (o *orders) Create(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error {
//		db := o.db
//		if txn != nil {
//			db = txn
//		}
//		return db.Create(order).Error
//	}
func (o *orders) Create(ctx context.Context, txn *sql.Tx, order *do.OrderInfoDO) error {
	_, err := txn.Exec("INSERT INTO orderinfo SET user=?, order_sn=?,order_mount=?,address=?,signer_name=?,singer_mobile=?,post=?",
		order.User, order.OrderSn, order.OrderMount, order.Address, order.SignerName, order.SingerMobile, order.Post)
	return err
}

func (o *orders) Update(ctx context.Context, txn *gorm.DB, order *do.OrderInfoDO) error {
	db := o.db
	if txn != nil {
		db = txn
	}
	return db.Model(order).Save(order).Error
}

var _ v1.OrderStore = &orders{}
