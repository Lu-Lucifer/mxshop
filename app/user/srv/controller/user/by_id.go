package user

import (
	"context"
	"fmt"
	upbv1 "mxshop/api/user/v1"
	"mxshop/pkg/log"
)

func (u *userServer) GetUserById(ctx context.Context, info *upbv1.IdRequest) (*upbv1.UserInfoResponse, error) {
	log.Infof("get user by id function is called")
	fmt.Println("srv GetUserById is called")
	user, err := u.srv.GetByID(ctx, uint64(info.Id))
	if err != nil {
		log.Errorf("get user by id:%d,error:%v", info.Id, err)
		return nil, err
	}
	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}
