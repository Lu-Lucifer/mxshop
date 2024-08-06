package user

import (
	"github.com/gin-gonic/gin"
	gin2 "mxshop/app/pkg/translator/gin"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/pkg/common/core"
	jtime "mxshop/pkg/common/time"
	"time"
)

type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,min=3,max=10"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-06"` // 生日格式 2006-01-02
}

func (us *userServer) UpdateUser(ctx *gin.Context) {
	updateForm := &UpdateUserForm{}
	if err := ctx.ShouldBind(updateForm); err != nil {
		gin2.HandleValidatorError(ctx, err, us.trans)
		return
	}

	//先通过userid查找用户
	userID, _ := ctx.Get(middlewares.KeyUserID)
	userIDInt := uint64(userID.(float64))
	userDTO, err := us.sf.User().Get(ctx, userIDInt)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}

	//更新用户
	loc, _ := time.LoadLocation("Local")
	birthday, _ := time.ParseInLocation("2006-01-02", updateForm.Birthday, loc)
	userDTO.NickName = updateForm.Name
	userDTO.Gender = updateForm.Gender
	userDTO.Birthday = jtime.Time{birthday}
	//TODO 修改时add_time插入空值报错
	err = us.sf.User().Update(ctx, userDTO)
	if err != nil {
		core.WriteResponse(ctx, err, nil)
		return
	}
	core.WriteResponse(ctx, nil, nil)

}
