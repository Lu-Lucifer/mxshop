package options

import (
	"github.com/spf13/pflag"
)

type EsOptions struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

func NewEsOptions() *EsOptions {
	return &EsOptions{
		Host: "127.0.0.1",
		Port: 9200,
	}
}

func (e *EsOptions) Validate() []error {
	var errs []error

	return errs
}

func (e *EsOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&e.Host, "es.host", e.Host, "elasticsearch host")
	fs.IntVar(&e.Port, "es.port", e.Port, "elasticsearch port")
}
