package v1

import (
	"context"
	"github.com/go-redsync/redsync/v4"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	v1 "mxshop/app/inventory/srv/internal/data/v1"
	"mxshop/app/inventory/srv/internal/domain/do"
	"mxshop/app/inventory/srv/internal/domain/dto"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	"mxshop/pkg/errors"
	"mxshop/pkg/log"
	"sort"
)

const (
	inventoryLockPrefix = "inventory_"
	orderLockPrefix     = "order_"
)

type InventorySrv interface {
	//设置库存
	Create(ctx context.Context, inv *dto.InventoryDTO) error
	//根据商品id查询库存
	Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error)
	//扣减库存
	Sell(ctx context.Context, ordersn string, detail []do.GoodsDetail) error
	//归还库存
	Reback(ctx context.Context, ordersn string, detail []do.GoodsDetail) error
}

type inventoryService struct {
	data         v1.DataFactory
	redisOptions *options.RedisOptions
	pool         redsyncredis.Pool
}

var _ InventorySrv = &inventoryService{}

func newInventoryService(s *service) *inventoryService {
	return &inventoryService{
		data:         s.data,
		redisOptions: s.redisOptions,
		pool:         s.pool,
	}
}
func (is *inventoryService) Create(ctx context.Context, inv *dto.InventoryDTO) error {
	return is.data.Inventorys().Create(ctx, &inv.InventoryDO)
}

func (is *inventoryService) Get(ctx context.Context, goodsID uint64) (*dto.InventoryDTO, error) {
	inv, err := is.data.Inventorys().Get(ctx, goodsID)
	if err != nil {
		return nil, err
	}
	return &dto.InventoryDTO{
		InventoryDO: *inv,
	}, nil
}

func (is *inventoryService) Sell(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s扣减库存", ordersn)
	rs := redsync.New(is.pool)
	//实际上批量扣减库存的时候，我们经常会先按照商品的id排序，然后从大到小逐个扣减库存，这样可以减少锁的竞争
	//先details类型转换为GoodsDetailList别名类型
	var detail = do.GoodsDetailList(details)
	sort.Sort(detail)

	txn := is.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("inventoryService.Sell事务进行中出现异常: %v", err)
			return
		}
	}()
	sellDetail := do.StockSellDetailDO{
		OrderSn: ordersn,
		Status:  1,
		Detail:  detail,
	}
	for _, goodsInfo := range detail {
		mutex := rs.NewMutex(inventoryLockPrefix + ordersn)
		if err := mutex.Lock(); err != nil {
			log.Errorf("订单%s获取锁失败", ordersn)
			return err
		}
		inv, err := is.data.Inventorys().Get(ctx, uint64(goodsInfo.Goods))
		if err != nil {
			log.Errorf("订单%s获取库存失败", ordersn)
			return err
		}
		//判断库存是否充足
		if inv.Stocks < goodsInfo.Num {
			txn.Rollback()
			log.Errorf("商品%d库存不足，现有库存：%d", goodsInfo.Goods, inv.Stocks)
			return errors.WithCode(code.ErrInvNotEnough, "库存不足")
		}
		//扣减库存
		inv.Stocks -= goodsInfo.Num
		err = is.data.Inventorys().Reduce(ctx, txn, uint64(goodsInfo.Goods), int(goodsInfo.Num))
		if err != nil {
			txn.Rollback()
			log.Errorf("订单%s扣减库存失败", ordersn)
			return err
		}
		//释放锁
		if _, err = mutex.Unlock(); err != nil {
			txn.Rollback()
			log.Errorf("订单%s释放锁失败", ordersn)
		}
	}

	err := is.data.Inventorys().CreateStockSellDetail(ctx, txn, &sellDetail)
	if err != nil {
		txn.Rollback()
		log.Errorf("订单%s创建扣减库存记录失败", ordersn)
		return err
	}

	txn.Commit()
	return nil
}

func (is *inventoryService) Reback(ctx context.Context, ordersn string, details []do.GoodsDetail) error {
	log.Infof("订单%s归还库存", ordersn)
	rs := redsync.New(is.pool)
	txn := is.data.Begin()
	defer func() {
		if err := recover(); err != nil {
			txn.Rollback()
			log.Errorf("inventoryService.Sell事务进行中出现异常: %v", err)
			return
		}
	}()
	//库存归还的时候有细节
	//1.主动取消 2.网络问题引起的重试 3.超时取消 4.退款取消   加锁保证：获取扣减库存记录
	mutex := rs.NewMutex(orderLockPrefix + ordersn)
	if err := mutex.Lock(); err != nil {
		txn.Rollback()
		log.Errorf("订单%s获取锁失败", ordersn)
		return err
	}

	sellDetail, err := is.data.Inventorys().GetSellDetail(ctx, txn, ordersn)
	if err != nil {
		txn.Rollback()
		_, err := mutex.Unlock()
		if err != nil {
			return err
		}
		if errors.IsCode(err, code.ErrInvSellDetailNotFound) {
			log.Errorf("订单%s扣减库存记录不存在，忽略", ordersn)
			return nil
		}
		log.Errorf("订单%s获取扣减库存记录失败", ordersn)
		return err
	}

	if sellDetail.Status == 2 {
		log.Infof("订单%s扣减库存记录已经归还，忽略", ordersn)
		return nil
	}

	var detail = do.GoodsDetailList(details)
	sort.Sort(detail)

	for _, goodsInfo := range detail {
		inv, err := is.data.Inventorys().Get(ctx, uint64(goodsInfo.Goods))
		if err != nil {
			txn.Rollback()
			log.Errorf("订单%s获取库存失败", ordersn)
			return err
		}

		//归还库存
		inv.Stocks += goodsInfo.Num
		err = is.data.Inventorys().Increase(ctx, txn, uint64(goodsInfo.Goods), int(goodsInfo.Num))
		if err != nil {
			txn.Rollback()
			log.Errorf("订单%s增加库存失败", ordersn)
			return err
		}
	}

	err = is.data.Inventorys().UpdateStockSellDetailStatus(ctx, txn, ordersn, 2)
	if err != nil {
		txn.Rollback()
		log.Errorf("订单%s更新扣减库存记录失败", ordersn)
		return err
	}

	txn.Commit()
	return nil
}
