package user

import (
	"context"

	upbv1 "mxshop/api/user/v1"
	v12 "mxshop/app/user/srv/data/v1"
	v1 "mxshop/app/user/srv/service/v1"
	"mxshop/pkg/log"
)

// controller层应该是很薄的一层，参数校验，日志打印，错误处理，调用service层
func (u *userServer) CreateUser(ctx context.Context, request *upbv1.CreateUserInfo) (*upbv1.UserInfoResponse, error) {
	log.Info("get user by mobile function called.")

	userDO := v12.UserDO{
		Mobile:   request.Mobile,
		NickName: request.NickName,
		Password: request.PassWord,
	}
	userDTO := v1.UserDTO{userDO}

	err := u.srv.Create(ctx, &userDTO)
	if err != nil {
		log.Errorf("get user by mobile: %s,error: %v", request.Mobile, err)
	}
	userInfoRsp := DTOToResponse(userDTO)
	return userInfoRsp, nil
}
