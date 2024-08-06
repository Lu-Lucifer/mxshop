package user

//包名这里不要用v1版本号，后期加版本不好维护
import (
	ut "github.com/go-playground/universal-translator"
	"mxshop/app/mxshop/api/internal/service"
)

type userServer struct {
	trans ut.Translator
	sf    service.ServiceFactory
}

func NewUserController(trans ut.Translator, sf service.ServiceFactory) *userServer {
	return &userServer{
		trans: trans,
		sf:    sf,
	}
}
