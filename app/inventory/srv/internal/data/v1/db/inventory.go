package mysql

import (
	"context"
	v1 "mxshop/app/inventory/srv/internal/data/v1"
	"mxshop/app/inventory/srv/internal/domain/do"
	"mxshop/app/pkg/code"
	code2 "mxshop/gmicro/code"
	"mxshop/pkg/errors"

	"gorm.io/gorm"
	"mxshop/pkg/log"
)

type inventorys struct {
	db *gorm.DB
}

// 更新库存销售状态
func (i *inventorys) UpdateStockSellDetailStatus(ctx context.Context, txn *gorm.DB, ordersn string, status int32) error {
	db := i.db
	if txn != nil {
		db = txn
	}

	//update语句如果没有更新的话那么不会报错，但是他会返回一个影响的行数，所以我们可以根据影响的行数来判断是否更新成功   原子操作
	result := db.Model(do.StockSellDetailDO{}).Where("order_sn = ?", ordersn).Update("status", status)
	if result.Error != nil {
		return errors.WithCode(code2.ErrDatabase, result.Error.Error())
	}

	//这里应该在service层去写代码判断更合理
	//有两种情况都会导致影响的行数为0，一种是没有找到，一种是没有更新
	//if result.RowsAffected == 0 {
	//	return errors.WithCode(code.ErrInvSellDetailNotFound, "inventory sell detail not found")
	//}
	return nil
}

// 查询库存销售信息
func (i *inventorys) GetSellDetail(ctx context.Context, txn *gorm.DB, ordersn string) (*do.StockSellDetailDO, error) {
	db := i.db
	if txn != nil {
		db = txn
	}
	var ordersellDetail do.StockSellDetailDO
	err := db.Where("order_sn = ?", ordersn).First(&ordersellDetail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrInvSellDetailNotFound, err.Error())
		}
		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}
	return &ordersellDetail, err
}

// 扣减库存
func (i *inventorys) Reduce(ctx context.Context, txn *gorm.DB, goodsID uint64, num int) error {
	db := i.db
	if txn != nil {
		db = txn
	}
	//采用了数据库的原子性功能   乐观锁
	return db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).Where("stocks >= ?", num).UpdateColumn("stocks", gorm.Expr("stocks - ?", num)).Error
}

// 新增库存
func (i *inventorys) Increase(ctx context.Context, txn *gorm.DB, goodsID uint64, num int) error {
	db := i.db
	if txn != nil {
		db = txn
	}
	err := db.Model(&do.InventoryDO{}).Where("goods=?", goodsID).UpdateColumn("stocks", gorm.Expr("stocks + ?", num)).Error
	return err
}

// 新增库存销售信息  扣减库存时调用  controller层不需要知道这些细节
func (i *inventorys) CreateStockSellDetail(ctx context.Context, txn *gorm.DB, detail *do.StockSellDetailDO) error {
	db := i.db
	if txn != nil {
		db = txn
	}

	tx := db.Create(&detail)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

// 新建库存信息
func (i *inventorys) Create(ctx context.Context, inv *do.InventoryDO) error {
	//设置库存， 如果我要更新库存
	tx := i.db.Create(&inv)
	if tx.Error != nil {
		return errors.WithCode(code2.ErrDatabase, tx.Error.Error())
	}
	return nil
}

// 查询商品的库存信息
func (i *inventorys) Get(ctx context.Context, goodsID uint64) (*do.InventoryDO, error) {
	inv := do.InventoryDO{}
	err := i.db.Where("goods = ?", goodsID).First(&inv).Error
	if err != nil {
		log.Errorf("get inv err: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(code.ErrInventoryNotFound, err.Error())
		}

		return nil, errors.WithCode(code2.ErrDatabase, err.Error())
	}

	return &inv, nil
}

func newInventorys(data *mysqlStore) *inventorys {
	return &inventorys{db: data.db}
}

var _ v1.InventoryStore = &inventorys{}
