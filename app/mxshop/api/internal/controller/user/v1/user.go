package user

//包名这里不要用v1版本号，后期加版本不好维护
import (
	ut "github.com/go-playground/universal-translator"
	"mxshop/app/mxshop/api/internal/service/user/v1"
)

type userServer struct {
	trans ut.Translator
	srv   user.UserSrv
}

func NewUserController(trans ut.Translator, srv user.UserSrv) *userServer {
	return &userServer{
		trans: trans,
		srv:   srv,
	}
}
