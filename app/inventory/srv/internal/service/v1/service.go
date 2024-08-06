package v1

import (
	"fmt"
	redsyncredis "github.com/go-redsync/redsync/v4/redis"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
	v1 "mxshop/app/inventory/srv/internal/data/v1"
	"mxshop/app/pkg/options"
)

type ServiceFactory interface {
	Inventorys() InventorySrv
}

type service struct {
	data         v1.DataFactory
	redisOptions *options.RedisOptions
	pool         redsyncredis.Pool
}

func (s *service) Inventorys() InventorySrv {
	return newInventoryService(s)
}

var _ ServiceFactory = &service{}

func NewService(data v1.DataFactory, redisOptions *options.RedisOptions) *service {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", redisOptions.Host, redisOptions.Port),
	})
	pool := goredis.NewPool(client)
	return &service{
		data:         data,
		redisOptions: redisOptions,
		pool:         pool,
	}
}
