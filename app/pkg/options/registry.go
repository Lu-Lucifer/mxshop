package options

import (
	"github.com/spf13/pflag"
	"mxshop/pkg/errors"
)

type RegistryOptions struct {
	Address string `mapstructure:"address" json:"address,omitempty"`
	Scheme  string `mapstructure:"scheme" json:"scheme,omitempty"`
}

func NewRegistryOptions() *RegistryOptions {
	return &RegistryOptions{
		Address: "127.0.0.1:8500",
		Scheme:  "http",
	}
}

func (ro *RegistryOptions) Validate() []error {
	errs := []error{}
	if ro.Address == "" || ro.Scheme == "" {
		errs = append(errs, errors.New("address and schema is empty"))
	}
	return errs
}

// 添加到命令组，可以从命令行读取服务注册的参数
func (ro *RegistryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&ro.Address, "consul.address", ro.Address, "consul address,default is 127.0.0.1:8500")
	fs.StringVar(&ro.Scheme, "consul.schema", ro.Scheme, "consul schema,default is http")
}
