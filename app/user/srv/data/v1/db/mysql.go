package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	"mxshop/pkg/errors"
	"sync"
)

/*
这个方法应该返回的是全局的一个变量
另外如果开始没有初始化好，那么就初始化一次，后续直接拿到db变量
*/

var (
	dbFactory *gorm.DB
	once      sync.Once
)

func GetDBFactoryOr(mysqlOpts *options.MySQLOptions) (*gorm.DB, error) {
	//校验一下，没有db实例，也不传配置
	if mysqlOpts == nil && dbFactory == nil {
		return nil, fmt.Errorf("failed to get mysql store factory")
	}
	var err error
	once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlOpts.Username,
			mysqlOpts.Password,
			mysqlOpts.Host,
			mysqlOpts.Port,
			mysqlOpts.Database,
		)

		dbFactory, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return
		}
		//设置线程池
		sqlDB, _ := dbFactory.DB()
		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(mysqlOpts.MaxIdleConnections)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(mysqlOpts.MaxOpenConnections)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(mysqlOpts.MaxConnectionLifetime)
	})
	if dbFactory == nil || err != nil {
		return nil, errors.WithCode(code.ErrConnectDB, "failed to get mysql factory")
	}
	return dbFactory, nil
}
