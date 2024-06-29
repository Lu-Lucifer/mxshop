package v1

import (
	"context"
	"mxshop/app/pkg/code"
	dv1 "mxshop/app/user/srv/data/v1"
	metav1 "mxshop/pkg/common/meta/v1"
	"mxshop/pkg/errors"
)

type UserDTO struct {
	dv1.UserDO
}
type UserDTOList struct {
	TotalCount int64      //总数
	Items      []*UserDTO //数据
}

// service层报漏这个接口给controller层调用
type UserSrv interface {
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	GetByID(ctx context.Context, id uint64) (*UserDTO, error)
	Create(ctx context.Context, user *UserDTO) error
	Update(ctx context.Context, user *UserDTO) error
}

var _ UserSrv = &userService{}

type userService struct {
	userStore dv1.UserStore
}

func (us *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	userDO, err := us.userStore.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return &UserDTO{UserDO: *userDO}, nil
}

func (us *userService) GetByID(ctx context.Context, id uint64) (*UserDTO, error) {
	userDO, err := us.userStore.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserDTO{UserDO: *userDO}, nil
}

func (us *userService) Create(ctx context.Context, user *UserDTO) error {
	//一个手机号码只能创建一个用户，且只能通过手机号码创建。判断用户是否存在
	_, err := us.userStore.GetByMobile(ctx, user.Mobile)
	//只有手机号不存在的情况下才能注册
	if err != nil && errors.IsCode(err, code.ErrUserNotFound) {
		return us.userStore.Create(ctx, &user.UserDO)
	}
	//这里应该区别到底是什么错误
	return errors.WithCode(code.ErrUserAlreadyExists, "用户已经存在")

}

func (us *userService) Update(ctx context.Context, user *UserDTO) error {
	//先查询用户是否存在 也可以不用查询
	_, err := us.userStore.GetByID(ctx, uint64(user.ID))
	if err != nil {
		return err
	}
	return us.userStore.Update(ctx, &user.UserDO)
}

func NewUserService(userStore dv1.UserStore) *userService {
	return &userService{
		userStore: userStore,
	}
}

func (us *userService) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error) {
	doList, err := us.userStore.List(ctx, orderby, opts)
	if err != nil {
		return nil, err
	}
	var userDTOList UserDTOList
	for _, value := range doList.Items {
		projectDTO := UserDTO{
			UserDO: *value,
		}
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}
	return &userDTOList, nil
}

//func List(ctx context.Context,opts metav1.ListMeta)(*UserDTOList,error){
//	doList, err := dv1.List(ctx,opts)
//	if err != nil {
//		return nil, err
//	}
//	var userDTOList UserDTOList
//	for _, value := range doList.Items {
//		projectDTO := UserDTO{
//            Name: value.Name,
//        }
//        userDTOList.Items = append(userDTOList.Items,&projectDTO)
//	}
//	return &userDTOList,nil
//}
