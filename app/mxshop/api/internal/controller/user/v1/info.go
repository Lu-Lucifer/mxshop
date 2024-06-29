package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (us *userServer) GetUserDetail(ctx *gin.Context) {
	//userid, _ := ctx.Get(middlewares.KeyUserID)
	//ctx.JSON(http.StatusOK, gin.H{
	//	"id": userid,
	//})
	userIDStr := ctx.Param("userid")
	atoi, _ := strconv.Atoi(userIDStr)
	userDTO, err := us.srv.Get(ctx, uint64(atoi))
	if err != nil {
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": userDTO,
	})
}
