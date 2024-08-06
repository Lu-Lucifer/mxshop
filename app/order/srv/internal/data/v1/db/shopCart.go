package db

import (
	"context"
	"database/sql"
	"gorm.io/gorm"
	v1 "mxshop/app/order/srv/internal/data/v1"
	"mxshop/app/order/srv/internal/domain/do"
	"mxshop/app/pkg/code"
	code2 "mxshop/gmicro/code"
	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/errors"
)

type shopCarts struct {
	db *gorm.DB
}

func newshopCarts(factory *dataFactory) *shopCarts {
	return &shopCarts{db: factory.db}
}
func (sc *shopCarts) List(ctx context.Context, userID uint64, checked bool, meta metav1.ListMeta, orderby []string) (*do.ShoppingCartDOList, error) {
	ret := &do.ShoppingCartDOList{}
	query := sc.db
	var limit, offset int
	if meta.PageSize == 0 {
		limit = 10
	} else {
		limit = meta.PageSize
	}
	if meta.Page > 0 {
		offset = (meta.Page - 1) * limit
	}
	if userID > 0 {
		query = query.Where("user = ?", userID)
	}
	if checked {
		query = query.Where("checked = ?", true)
	}
	//排序
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

func (sc *shopCarts) Create(ctx context.Context, cartItem *do.ShoppingCartDO) error {
	tx := sc.db.Create(cartItem)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

func (sc *shopCarts) Get(ctx context.Context, userID, goodsID uint64) (*do.ShoppingCartDO, error) {
	var shopCart do.ShoppingCartDO
	//gorm链路追踪
	err := sc.db.WithContext(ctx).Where("user = ? AND goods = ?", userID, goodsID).First(&shopCart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if errors.IsCode(err, code.ErrShopCartItemNotFound) {
				return nil, errors.WithCode(code2.ErrDatabase, err.Error())
			}
		}
	}
	return &shopCart, err
}

func (sc *shopCarts) UpdateNum(ctx context.Context, cartItem *do.ShoppingCartDO) error {
	return sc.db.Model(&do.ShoppingCartDO{}).Where("user = ? AND goods=?", cartItem.User, cartItem.Goods).Update("nums", cartItem.Nums).Update("checked", cartItem.Checked).Error
}

func (sc *shopCarts) Delete(ctx context.Context, ID uint64) error {
	return sc.db.Where("id = ?", ID).Delete(&do.ShoppingCartDO{}).Error
}

// 清空check状态
func (sc *shopCarts) ClearCheck(ctx context.Context, userID uint64) error {
	panic("")
}

/*
	删除选中商品的购物车记录，下订单了
	从架构上来讲，这种实现有两种方案
	下单后，直接执行删除购物车的记录，比较简单
	下单后什么都不做，直接给rocketmq发送一个消息，然后由rocketmq来执行删除购物车的记录
*/

// 这个在事务中执行，建议使用消息队列
//
//	func (sc *shopCarts) DeleteByGoodsIDs(ctx context.Context, txn *gorm.DB, userID uint64, goodsIDs []int32) error {
//		db := sc.db
//		if txn != nil {
//			db = txn
//		}
//		return db.Where("user = ? AND goods IN (?)", userID, goodsIDs).Delete(&do.ShoppingCartDO{}).Error
//	}
func (sc *shopCarts) DeleteByGoodsIDs(ctx context.Context, txn *sql.Tx, userID uint64, goodsIDs []int32) error {
	_, err := txn.Exec("delete from shoppingcart where user = ? and goods in (?,?)", userID, goodsIDs[0], goodsIDs[1])
	return err
}

var _ v1.ShopCartStore = &shopCarts{}
