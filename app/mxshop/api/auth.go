package api

import (
	ginjwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/gmicro/server/restserver/middlewares/auth"
)

func newJWTAuth(opts *options.JwtOptions) middlewares.AuthStrategy {
	gjwt, _ := ginjwt.New(&ginjwt.GinJWTMiddleware{
		Realm:            opts.Realm,
		SigningAlgorithm: "HS256",
		Key:              []byte(opts.Key),
		Timeout:          opts.Timeout,
		MaxRefresh:       opts.MaxRefresh,
		LogoutResponse: func(c *gin.Context, code int) {
			c.JSON(code, nil)
		},
		IdentityHandler: claimHandlerFunc,
		IdentityKey:     middlewares.KeyUserID,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
	})
	return auth.NewJWTStrategy(*gjwt)
}

func claimHandlerFunc(c *gin.Context) interface{} {
	claims := ginjwt.ExtractClaims(c)
	c.Set(middlewares.KeyUserID, claims[middlewares.KeyUserID])
	return claims[ginjwt.IdentityKey]
}
