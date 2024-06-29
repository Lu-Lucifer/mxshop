package options

import (
	"github.com/spf13/pflag"
	"mxshop/pkg/errors"
)

type TelemetryOptions struct {
	Name     string  `json:"name"`
	Endpoint string  `json:"endpoint"`
	Sampler  float64 `json:"sampler"` //采样率
	Batcher  string  `json:"batcher"`
}

func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Name:     "mxshop",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	}
}

func (to *TelemetryOptions) Validate() []error {
	errs := []error{}
	if to.Batcher != "jaeger" && to.Batcher != "zipkin" {
		errs = append(errs, errors.New("address and schema is empty"))
	}
	return errs
}

// 添加到命令组，可以从命令行读取服务注册的参数
func (to *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&to.Name, "telemetey.name", to.Name, "opentelemetry name")
	fs.StringVar(&to.Endpoint, "telemetey.endpoint", to.Endpoint, "opentelemetry endpoint")
	fs.Float64Var(&to.Sampler, "telemetey.sampler", to.Sampler, "opentelemetry sampler")
	fs.StringVar(&to.Batcher, "telemetey.batcher", to.Batcher, "opentelemetry batcher,only support jaeger and zipkin")
}
