package options

import (
	"github.com/spf13/pflag"
	"mxshop/pkg/errors"
)

type DtmOptions struct {
	GrpcServer string `json:"grpc" mapstructure:"grpc,omitempty"`
	HttpServer string `json:"http" mapstructure:"http,omitempty"`
}

func NewDtmOptionsOptions() *DtmOptions {
	return &DtmOptions{
		HttpServer: "http://127.0.0.1:36789/api/dtmsvr",
		GrpcServer: "127.0.0.1:36790",
	}
}
func (o *DtmOptions) Validate() []error {
	errs := []error{}
	if o.HttpServer == "" && o.GrpcServer == "" {
		errs = append(errs, errors.New("address and http/grpc server is empty"))
	}
	return errs
}

func (o *DtmOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.GrpcServer, "dtm.grpcserver", o.GrpcServer, ""+
		"DTM GrpcServer, if left , default is 127.0.0.1:36790")
	fs.StringVar(&o.HttpServer, "dtm.httpserver", o.HttpServer, ""+
		"DTM HttpServer, if left , default is http://127.0.0.1:36789/api/dtmsvr")

}
