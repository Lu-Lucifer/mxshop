package user

import (
	"github.com/gin-gonic/gin"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/pkg/log"
	"net/http"
)

type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` //自定义validator 验证手机号
	Password  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"` //图形验证码
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

func (us *userServer) Login(ctx *gin.Context) {
	log.Infof("Login is called")
	//表单验证

	passwordLoginForm := PasswordLoginForm{}
	// shouldbind自动识别form表单请求或json请求
	if err := ctx.ShouldBind(&passwordLoginForm); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	// 图形验证码校验
	ok := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true)
	//ok := true
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"captcha": "图形验证码错误",
		})
		return
	}

	userDTO, err := us.srv.MobileLogin(ctx, passwordLoginForm.Mobile, passwordLoginForm.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "登录失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":        userDTO.ID,
		"nick_name": userDTO.NickName,
		"token":     userDTO.Token,
		"expire_at": userDTO.ExpiresAt, //毫秒级别
	})

}
