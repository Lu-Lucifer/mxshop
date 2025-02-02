package user

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"mxshop/app/mxshop/api/internal/data"
	"mxshop/app/pkg/code"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/server/restserver/middlewares"
	"mxshop/pkg/errors"
	"mxshop/pkg/log"
	"mxshop/pkg/storage"
	"time"
)

type UserDTO struct {
	data.User
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type UserSrv interface {
	MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error)
	Register(ctx context.Context, mobile, password, codes string) (*UserDTO, error)
	Update(ctx context.Context, uerDTO *UserDTO) error
	Get(ctx context.Context, userID uint64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	CheckPassword(ctx context.Context, password, EncryptedPassword string) (bool, error)
}

type userService struct {
	ud      data.DataFactory
	jwtOpts *options.JwtOptions
}

func NewUserService(data data.DataFactory, jwtOpts *options.JwtOptions) UserSrv {
	return &userService{
		ud:      data,
		jwtOpts: jwtOpts,
	}
}

var _ UserSrv = &userService{}

func (u *userService) MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error) {
	user, err := u.ud.User().GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	//检查密码是否正确
	err = u.ud.User().CheckPassWord(ctx, password, user.PassWord)
	if err != nil {
		return nil, err
	}

	//生成token
	j := middlewares.NewJWT(u.jwtOpts.Key)
	token, err := j.CreateToken(middlewares.CustomClaims{
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                //生效时间
			ExpiresAt: time.Now().Local().Add(u.jwtOpts.Timeout).Unix(), //过期时间
			Issuer:    u.jwtOpts.Realm,
		},
	})
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		User:      *user,
		Token:     token,
		ExpiresAt: time.Now().Local().Add(u.jwtOpts.Timeout).Unix(),
	}, nil

}

func (u *userService) Register(ctx context.Context, mobile, password, codes string) (*UserDTO, error) {
	//短信验证码校验
	rstore := storage.RedisCluster{}
	value, err := rstore.GetKey(ctx, mobile)
	if err != nil {
		return nil, errors.WithCode(code.ErrCodeNotExist, "手机验证码不存在")
	}
	if value != codes {
		return nil, errors.WithCode(code.ErrCodeIncorrect, "手机验证码错误")
	}

	// 创建用户
	var user = &data.User{
		Mobile:   mobile,
		PassWord: password,
	}
	err = u.ud.User().Create(ctx, user)
	if err != nil {
		log.Errorf("user register failed: %v", err)
		return nil, err
	}
	//生成token
	j := middlewares.NewJWT(u.jwtOpts.Key)
	token, err := j.CreateToken(middlewares.CustomClaims{
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                //生效时间
			ExpiresAt: time.Now().Local().Add(u.jwtOpts.Timeout).Unix(), //过期时间
			Issuer:    u.jwtOpts.Realm,
		},
	})
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		User:      *user,
		Token:     token,
		ExpiresAt: time.Now().Local().Add(u.jwtOpts.Timeout).Unix(),
	}, nil

}

func (u *userService) Update(ctx context.Context, uerDTO *UserDTO) error {
	err := u.ud.User().Update(ctx, &data.User{
		ID:       uerDTO.ID,
		NickName: uerDTO.NickName,
		Birthday: uerDTO.Birthday,
		Gender:   uerDTO.Gender,
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *userService) Get(ctx context.Context, userID uint64) (*UserDTO, error) {
	user, err := u.ud.User().Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		User: *user,
	}, nil
}

func (u *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	user, err := u.ud.User().GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}
	return &UserDTO{
		User: *user,
	}, nil
}

func (u *userService) CheckPassword(ctx context.Context, password, EncryptedPassword string) (bool, error) {
	return false, nil
}
