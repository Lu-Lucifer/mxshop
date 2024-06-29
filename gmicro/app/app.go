package app

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"mxshop/gmicro/registry"
	gs "mxshop/gmicro/server"
	"mxshop/pkg/log"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type App struct {
	opts     options
	instance *registry.ServiceInstance
	lk       sync.Mutex
	cancel   func()
}

func New(opts ...Option) *App {
	//设置options的默认值
	o := options{
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10000 * time.Second,
		stopTimeout:      10000 * time.Second,
	}
	//设置ID默认值
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &App{opts: o}
}

// 启动整个服务
func (a *App) Run() error {
	//服务注册的信息
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}

	//将注册instance变成app字段，否则Stop方法拿不到这个实例。 这个变量可能被其他的goroutine访问到
	a.lk.Lock()
	a.instance = instance
	a.lk.Unlock()

	//启动rest server，监听端口
	//if a.opts.restServer != nil {
	//	go func() {
	//		err := a.opts.restServer.Start(context.Background())
	//		if err != nil {
	//			panic(any(err))
	//		}
	//	}()
	//}

	//启动rpc server，监听端口
	//if a.opts.rpcServer != nil {
	//	go func() {
	//		err := a.opts.rpcServer.Start(context.Background())
	//		if err != nil {
	//			panic(any(err))
	//		}
	//	}()
	//}

	/*
		现在启动了2个server
		这两个server是否必须同时启动成功？
		如果有一个启动失败，那么我们就要停止另外一个server
		如果启动了多个，如果其中一个启动失败，其他的应该被cancel
		使用errgroup来保证出错可以cancel掉
	*/

	//启动rest server，监听端口
	//if a.opts.restServer != nil {
	//	//监听启动rest server是否被cancel，如果cancel掉就stop服务
	//	eg.Go(func() error {
	//		<-ctx.Done()
	//		//为以防调用Stop方法超时，以至一直调用，传入context timeout
	//		sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
	//		defer cancel()
	//		return a.opts.restServer.Stop(sctx)
	//	})
	//	eg.Go(func() error {
	//		return a.opts.restServer.Start(ctx)
	//	})
	//}

	//if a.opts.rpcServer != nil {
	//	//监听启动rpc server是否被cancel，如果cancel掉就stop服务
	//	eg.Go(func() error {
	//		<-ctx.Done()
	//		//为以防调用Stop方法超时，以至无休止的等待，传入context timeout
	//		sctx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
	//		defer cancel()
	//		return a.opts.rpcServer.Stop(sctx)
	//	})
	//	eg.Go(func() error {
	//		return a.opts.rpcServer.Start(ctx)
	//	})
	//}

	//上面errgroup的重复代码太多，改造一下，抽象成一个Server接口
	srvs := []gs.Server{}
	if a.opts.restServer != nil {
		srvs = append(srvs, a.opts.restServer)
	}
	if a.opts.rpcServer != nil {
		srvs = append(srvs, a.opts.rpcServer)
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range srvs {
		//监听启动rest server是否被cancel，如果cancel掉就stop服务
		srv := srv //防止循环被替换server
		eg.Go(func() error {
			<-ctx.Done()
			//为以防调用Stop方法超时，以至一直调用，传入context timeout
			sctx, scancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer scancel()
			return srv.Stop(sctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			log.Info("start server")
			return srv.Start(ctx)
		})
	}
	//保证上面两个server成功，才能执行下面服务注册逻辑，使用waitgroup处理
	wg.Wait()

	//注册服务
	if a.opts.registrar != nil {
		//注册服务的时候，是网络连接，选择context timeout
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.registrarTimeout)
		defer rcancel()
		err = a.opts.registrar.Register(rctx, instance)
		if err != nil {
			log.Errorf("register service error:%s", err)
			return err
		}
	}

	//监听退出信息
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	//这里监听signal，退出app stop，也需要通知http server 的 stop
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			return a.Stop()
		}
	})
	//eg.Wait()会等待所有eg.Go执行结束
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// 注销服务
func (a *App) Stop() error {
	log.Info("Deregister service")
	a.lk.Lock()
	instance := a.instance
	a.lk.Unlock()
	if a.opts.registrar != nil && instance != nil {
		//注销服务
		rctx, rcancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
		defer rcancel()
		err := a.opts.registrar.Deregister(rctx, instance)
		if err != nil {
			log.Errorf("deregister service error:%s", err)
			return err
		}
	}
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

// 创建服务注册结构体
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0)
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if a.opts.rpcServer != nil {
		u := &url.URL{
			Scheme: "grpc",
			Host:   a.opts.rpcServer.Address(),
		}
		endpoints = append(endpoints, u.String())
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Endpoints: endpoints,
		//Version: "1.0",
	}, nil

}
