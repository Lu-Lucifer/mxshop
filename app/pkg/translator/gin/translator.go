package gin

import (
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

func HandleValidatorError(c *gin.Context, err error, trans ut.Translator) {
	// 获取validator.ValidationErrors类型的errors
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}
	// validator.ValidationErrors类型错误则进行翻译
	c.JSON(http.StatusOK, gin.H{
		"msg": removeTopStruct(errs.Translate(trans)),
	})
	return

}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
