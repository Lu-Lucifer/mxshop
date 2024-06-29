package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func pickNums() {

}

var Middlewares = defaultMiddlewares()

// 在map中储存中间件
func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery": gin.Recovery(),
		"cors":     Cors(),
		"context":  Context(),
		"trace":    otelgin.Middleware("my gin opentelemetry"),
	}
}
