package user

import (
	v1 "mxshop/api/user/v1"
	srv1 "mxshop/app/user/srv/service/v1"
)

type userServer struct {
	srv srv1.UserSrv
	v1.UnimplementedUserServer
}

func NewUserServer(srv srv1.UserSrv) *userServer {
	return &userServer{
		srv: srv,
	}
}
