package api

import (
	"mxshop/app/mxshop/api/config"
	v12 "mxshop/app/mxshop/api/internal/controller/sms/v1"
	user "mxshop/app/mxshop/api/internal/controller/user/v1"
	"mxshop/app/mxshop/api/internal/data/rpc"
	"mxshop/app/mxshop/api/internal/service/sms/v1"
	user2 "mxshop/app/mxshop/api/internal/service/user/v1"
	"mxshop/gmicro/server/restserver"
)

func initRouter(g *restserver.Server, cfg *config.Config) {
	v1 := g.Group("/v1")
	ugroup := v1.Group("/user")
	userData, err := rpc.GetDataFactoryOr(cfg.Registry)
	if err != nil {
		panic(any(err))
	}

	userService := user2.NewUserService(userData, cfg.Jwt)
	ucontroller := user.NewUserController(g.Translator(), userService)
	{
		ugroup.POST("pwd_login", ucontroller.Login)
		ugroup.POST("register", ucontroller.Register)
		ugroup.GET("getInfo/:userid", ucontroller.GetUserDetail)
	}

	baseRouter := v1.Group("/base")
	{
		smsService := sms.NewSmsService(cfg.Sms)
		smsCtl := v12.NewSmsController(smsService, g.Translator())
		baseRouter.POST("send_sms", smsCtl.SendSms)
		baseRouter.GET("captcha", user.GetCaptcha)
	}

}
