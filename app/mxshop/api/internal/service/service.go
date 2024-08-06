package service

import (
	"mxshop/app/mxshop/api/internal/data"
	v1 "mxshop/app/mxshop/api/internal/service/goods/v1"
	"mxshop/app/mxshop/api/internal/service/sms/v1"
	user2 "mxshop/app/mxshop/api/internal/service/user/v1"
	"mxshop/app/pkg/options"
)

type ServiceFactory interface {
	Goods() v1.GoodsSrv
	User() user2.UserSrv
	Sms() sms.SmsSrv
}

type service struct {
	data    data.DataFactory
	jwtOpts *options.JwtOptions
	smsOpts *options.SmsOptions
}

var _ ServiceFactory = &service{}

func (s *service) Sms() sms.SmsSrv {
	return sms.NewSmsService(s.smsOpts)
}
func (s *service) Goods() v1.GoodsSrv {
	return v1.NewGoods(s.data)
}

func (s *service) User() user2.UserSrv {
	return user2.NewUserService(s.data, s.jwtOpts)
}

func NewService(store data.DataFactory, jwtOpts *options.JwtOptions, smsOpts *options.SmsOptions) *service {
	return &service{
		data:    store,
		jwtOpts: jwtOpts,
		smsOpts: smsOpts,
	}
}
