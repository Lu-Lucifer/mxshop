package user

import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"net/http"
)

var store = base64Captcha.DefaultMemStore

func GetCaptcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80)
	captcha := base64Captcha.NewCaptcha(driver, store)
	id, b64s, err := captcha.Generate()
	if err != nil {
		zap.S().Errorf("生成图形验证码错误:%v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成图形验证码错误",
		})
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"captchaId": id,
		"picPath":   b64s,
	})

}
