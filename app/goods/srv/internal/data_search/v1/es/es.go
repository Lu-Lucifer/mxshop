package es

import (
	"github.com/olivere/elastic/v7"
	v1 "mxshop/app/goods/srv/internal/data_search/v1"
	"mxshop/app/pkg/db"
	"mxshop/app/pkg/options"
	"mxshop/pkg/errors"
	"sync"
)

var (
	//esClient *elastic.Client
	//工厂模式
	searchFactory v1.SearchFactory
	once          sync.Once
)

//func GetSeatchFactoryOr(opts *options.EsOptions) (*elastic.Client, error) {
//	if opts == nil && esClient == nil {
//		return nil, errors.New("failed to get as client")
//	}
//	var err error
//	once.Do(func() {
//		esOpt := db.EsOptions{
//			Host: opts.Host,
//			Port: opts.Port,
//		}
//		esClient, err = db.NewEsClient(&esOpt)
//		if err != nil {
//			return
//		}
//	})
//	if esClient == nil {
//		return nil, errors.New("failed to get es as client")
//	}
//	return esClient, err
//}

// 工厂模式
func GetSeatchFactoryOr(opts *options.EsOptions) (v1.SearchFactory, error) {
	if opts == nil && searchFactory == nil {
		return nil, errors.New("failed to get as client")
	}
	once.Do(func() {
		esOpt := db.EsOptions{
			Host: opts.Host,
			Port: opts.Port,
		}
		esCli, err := db.NewEsClient(&esOpt)
		if err != nil {
			return
		}
		//构造工厂实例
		searchFactory = &dataSearch{
			esClient: esCli,
		}
	})
	if searchFactory == nil {
		return nil, errors.New("failed to get es as client")
	}
	return searchFactory, nil
}

type dataSearch struct {
	esClient *elastic.Client
}

func (ds *dataSearch) Goods() v1.GoodsStore {
	return newGoods(ds)
}

var _ v1.SearchFactory = &dataSearch{}
