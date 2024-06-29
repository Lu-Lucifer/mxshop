package v1

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"mxshop/app/mxshop/api/internal/service/sms/v1"
	"mxshop/app/pkg/code"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/pkg/common/core"
	"mxshop/pkg/errors"
	"mxshop/pkg/storage"
	"time"
)

type SmsController struct {
	srv   sms.SmsSrv
	trans ut.Translator
}

func NewSmsController(srv sms.SmsSrv, trans ut.Translator) *SmsController {
	return &SmsController{srv, trans}
}

type SendSmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` //自定义validator 验证手机号
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`
}

func (sc *SmsController) SendSms(ctx *gin.Context) {
	sendSmsForm := SendSmsForm{}
	// shouldbind自动识别form表单请求或json请求
	if err := ctx.ShouldBind(&sendSmsForm); err != nil {
		gin2.HandleValidatorError(ctx, err, sc.trans)
		return
	}
	smsCode := sms.GenerateSmsCode(6)
	err := sc.srv.SendSms(ctx, sendSmsForm.Mobile, "SMS_154950909", "{\"code\":"+smsCode+"}")
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}
	//将手机验证码保存在redis中
	rstore := storage.RedisCluster{}
	err = rstore.SetKey(ctx, sendSmsForm.Mobile, smsCode, 5*time.Minute)
	if err != nil {
		core.WriteResponse(ctx, errors.WithCode(code.ErrSmsSend, err.Error()), nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)

}
