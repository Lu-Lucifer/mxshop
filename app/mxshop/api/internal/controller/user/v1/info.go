package user

import (
	"github.com/gin-gonic/gin"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/pkg/common/core"
)

func (us *userServer) GetUserDetail(ctx *gin.Context) {
	//获取jwt中间件中设置的userid
	userID, _ := ctx.Get(middlewares.KeyUserID)
	userDTO, err := us.srv.Get(ctx, uint64(userID.(float64)))
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, err, gin.H{
		"name":     userDTO.NickName,
		"birthday": userDTO.Birthday.Format("2006-01-02"),
		"gender":   userDTO.Gender,
		"mobile":   userDTO.Mobile,
	})
}
