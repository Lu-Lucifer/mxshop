package app

import (
	"mxshop/gmicro/registry"
	"mxshop/gmicro/server/restserver"
	"mxshop/gmicro/server/rpcserver"
	"net/url"
	"os"
	"time"
)

type options struct {
	endpoints []*url.URL
	id        string
	name      string
	sigs      []os.Signal
	//允许用户传入自己的服务实现
	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
	//传递rpc服务
	rpcServer *rpcserver.Server
	//传递rest服务
	restServer *restserver.Server
}
type Option func(o *options)

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}

func WithRPCServer(rpcServer *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = rpcServer
	}
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithRestServer(restServer *restserver.Server) Option {
	return func(o *options) {
		o.restServer = restServer
	}
}
