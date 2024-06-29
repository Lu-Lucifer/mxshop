package restserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	mws "mxshop/gmicro/server/restserver/middlewares"
	"mxshop/gmicro/server/restserver/pprof"
	"mxshop/gmicro/server/restserver/validation"
	"mxshop/pkg/errors"
	"mxshop/pkg/log"
	"net/http"
	"time"
)

type JwtInfo struct {
	Realm      string
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
}

type Server struct {
	*gin.Engine
	port int
	//开发模式，默认debug
	mode string
	//是否开启健康检查，默认开启，如果开启会自动添加/health接口
	healthz bool
	//是否开启pprof,默认开启，如果开启会自动添加/debug/pprof接口
	enableProfiling bool
	//中间件
	middlewares []string
	jwt         *JwtInfo
	//翻译器语言
	transName string
	trans     ut.Translator
	server    *http.Server
}

func NewServer(opts ...SeverOption) *Server {
	srv := &Server{
		port:            8080,
		mode:            "debug",
		healthz:         true,
		enableProfiling: true,
		jwt: &JwtInfo{
			Realm:      "JWT",
			Key:        "^RJ34rXWbrrV@96#qXd3CaI%qpXEwV%#",
			Timeout:    time.Hour * 7 * 24,
			MaxRefresh: time.Hour * 7 * 24,
		},
		Engine:    gin.Default(),
		transName: "zh",
	}
	for _, opt := range opts {
		opt(srv)
	}
	//加载中间件
	for _, m := range srv.middlewares {
		mw, ok := mws.Middlewares[m]
		if !ok {
			log.Warnf("can not find middleware:%s", m)
			continue
		}
		srv.Use(mw)
	}
	return srv
}
func (s *Server) Translator() ut.Translator {
	return s.trans
}
func (s *Server) Start(ctx context.Context) error {
	if s.mode != gin.DebugMode && s.mode != gin.ReleaseMode && s.mode != gin.TestMode {
		return errors.New("mode must be one of debug/release/test")
	}
	//设置开发模式，打印路由信息
	gin.SetMode(s.mode)
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Infof("%-6s %-s --> %s(%d handlers)", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	err := s.InitTrans(s.transName)
	if err != nil {
		log.Errorf("initTrans error: %s", err.Error())
		return err
	}

	//注册mobile验证器
	validation.RegisterMobile(s.trans)

	//根据配置初始化pprof
	if s.enableProfiling {
		pprof.Register(s.Engine)
	}

	log.Infof("rest server is running on %d", s.port)
	address := fmt.Sprintf(":%d", s.port)
	// 用http的sever来启动服务
	s.server = &http.Server{
		Addr:    address,
		Handler: s.Engine,
	}
	err = s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Infof("rest server is stopping")
	//用http sever来关闭服务
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Errorf("rest server shutdown error: %s", err.Error())
		return err
	}
	log.Infof("rest server stopped")
	return nil
}
