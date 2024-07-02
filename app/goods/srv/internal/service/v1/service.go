package v1

import (
	v1 "mxshop/app/goods/srv/internal/data/v1"
	v12 "mxshop/app/goods/srv/internal/data_search/v1"
)

type ServiceFactory interface {
	Goods() GoodsSrv
}

type service struct {
	data       v1.DataFactory
	dataSearch v12.SearchFactory
}

func NewService(data v1.DataFactory, dataSearch v12.SearchFactory) *service {
	return &service{
		data:       data,
		dataSearch: dataSearch,
	}
}

func (s *service) Goods() GoodsSrv {
	return newGoods(s)
}

var _ ServiceFactory = &service{}
