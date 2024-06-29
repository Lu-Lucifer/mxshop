package config

import (
	"mxshop/app/pkg/options"
	"mxshop/pkg/app"
	cliflag "mxshop/pkg/common/cli/flag"
	"mxshop/pkg/log"
)

// 自定义config（从配置文件读取或从命令行获取）
type Config struct {
	Log       *log.Options              `mapstructure:"log" json:"log"`
	Server    *options.ServerOptions    `mapstructure:"server" json:"server"`
	Registry  *options.RegistryOptions  `mapstructure:"registry" json:"registry"`
	Telemetry *options.TelemetryOptions `mapstructure:"telemetry" json:"telemetry"`
}

func (c *Config) Flags() (fss cliflag.NamedFlagSets) {
	c.Log.AddFlags(fss.FlagSet("logs"))
	c.Server.AddFlags(fss.FlagSet("server"))
	c.Registry.AddFlags(fss.FlagSet("registry"))
	c.Registry.AddFlags(fss.FlagSet("telemetry"))
	return fss
}
func (c *Config) Validate() []error {
	var errs []error
	errs = append(errs, c.Log.Validate()...)
	errs = append(errs, c.Server.Validate()...)
	errs = append(errs, c.Registry.Validate()...)
	errs = append(errs, c.Telemetry.Validate()...)
	return errs
}

// 实现app.CliOptions接口
var _ app.CliOptions = &Config{}

func New() *Config {
	//配置默认初始化
	return &Config{
		Log:       log.NewOptions(),
		Server:    options.NewServerOptions(),
		Registry:  options.NewRegistryOptions(),
		Telemetry: options.NewTelemetryOptions(),
	}
}
