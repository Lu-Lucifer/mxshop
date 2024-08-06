package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	proto2 "mxshop/api/goods/v1"
	proto "mxshop/api/inventory/v1"
	v1 "mxshop/app/order/srv/internal/data/v1"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	"mxshop/pkg/errors"
	"os"
	"sync"
	"time"
)

type dataFactory struct {
	db          *gorm.DB
	invClient   proto.InventoryClient
	goodsClient proto2.GoodsClient
}

func (df *dataFactory) Orders() v1.OrderStore {
	return newOrders(df)
}

func (df *dataFactory) ShopCarts() v1.ShopCartStore {
	return newshopCarts(df)
}

func (df *dataFactory) Goods() proto2.GoodsClient {
	return df.goodsClient
}

func (df *dataFactory) Inventorys() proto.InventoryClient {
	return df.invClient
}

func (df *dataFactory) Begin() *gorm.DB {
	return df.db.Begin()
}

var _ v1.DataFactory = &dataFactory{}

var (
	data v1.DataFactory
	once sync.Once
)

func GetDataFactoryOr(mysqlOpts *options.MySQLOptions, registry *options.RegistryOptions) (v1.DataFactory, error) {
	if (mysqlOpts == nil && registry == nil) && data == nil {
		return nil, errors.New("failed to get data store factory")
	}
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database)
		//可以自己封装log
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				//慢查询阈值可以自己设置
				SlowThreshold:             time.Second,                         // Slow SQL threshold
				LogLevel:                  logger.LogLevel(mysqlOpts.LogLevel), // Log level
				IgnoreRecordNotFoundError: true,                                // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,                                // Don't include params in the SQL log
				Colorful:                  false,                               // Disable color
			},
		)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return
		}
		sqlDB, _ := db.DB()

		//允许连接多少个mysql
		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)
		//允许最大的空闲的连接数
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)
		//重用连接的最大时长
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)

		//服务发现
		goodsClient := GetGoodsClient(registry)
		invClient := GetInventoryClient(registry)

		data = &dataFactory{
			db:          db,
			goodsClient: goodsClient,
			invClient:   invClient,
		}
	})
	var err error
	if data == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get data store factory")
	}
	return data, nil
}
