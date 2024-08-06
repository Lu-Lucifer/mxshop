package user

import (
	"github.com/gin-gonic/gin"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/pkg/common/core"
)

type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"` //自定义validator 验证手机号
	Password string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"` //短信验证码

}

func (us *userServer) Register(ctx *gin.Context) {
	registerForm := RegisterForm{}
	if err := ctx.ShouldBind(&registerForm); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	//短信验证码校验
	userDTO, err := us.sf.User().Register(ctx, registerForm.Mobile, registerForm.Password, registerForm.Code)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	core.WriteResponse(ctx, nil, gin.H{
		"id":        userDTO.ID,
		"nick_name": userDTO.NickName,
		"token":     userDTO.Token,
		"expire_at": userDTO.ExpiresAt,
	})

}
