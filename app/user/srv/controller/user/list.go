package user

import (
	"context"
	upbv1 "mxshop/api/user/v1"
	srvv1 "mxshop/app/user/srv/service/v1"
	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/log"
)

func DTOToResponse(userDTO srvv1.UserDTO) *upbv1.UserInfoResponse {
	userInfoRsp := upbv1.UserInfoResponse{
		Id:       userDTO.ID,
		PassWord: userDTO.Password,
		Mobile:   userDTO.Mobile,
		NickName: userDTO.NickName,
		//在grpc中message字段有默认值,不能随便赋值nil进去，容易出错
		// BirthDay: user.Birthday,
		Gender: userDTO.Gender,
		Role:   int32(userDTO.Role),
	}
	if userDTO.Birthday != nil {
		// 时间类型转换为int类型
		userInfoRsp.BirthDay = uint64(userDTO.Birthday.Unix())
	}

	return &userInfoRsp
}

//func GetUserList(ctx context.Context,info *upbv1.PageInfo)(*upbv1.UserListResponse, error){
//	srvOpts := metav1.ListMeta{
//		Page: int(info.Pn),
//		PageSize: int(info.PSize),
//	}
//	dtoList, err := srvv1.List(ctx,srvOpts)
//	if err != nil {
//		return nil, err
//	}
//	var rsp upbv1.UserListResponse
//	for _, value := range dtoList.Items {
//		userRsp := DTOToResponse(*value)
//		rsp.Data = append(rsp.Data,&userRsp)
//	}
//	return &rsp,nil
//}

func (u *userServer) GetUserList(ctx context.Context, info *upbv1.PageInfo) (*upbv1.UserListResponse, error) {

	log.Info("GetUserList is called")
	srvOpts := metav1.ListMeta{
		Page:     int(info.Pn),
		PageSize: int(info.PSize),
	}

	dtoList, err := u.srv.List(ctx, []string{}, srvOpts)
	if err != nil {
		return nil, err
	}
	var rsp = upbv1.UserListResponse{Total: 23}
	for _, value := range dtoList.Items {
		userRsp := DTOToResponse(*value)
		rsp.Data = append(rsp.Data, userRsp)
	}
	return &rsp, nil
}
