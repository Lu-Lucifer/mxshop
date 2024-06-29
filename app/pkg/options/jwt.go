package options

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/spf13/pflag"
	"time"
)

type JwtOptions struct {
	Realm      string        `mapstructure:"realm" json:"realm"`
	Key        string        `mapstructure:"key" json:"key"`
	Timeout    time.Duration `mapstructure:"timeout" json:"timeout"`
	MaxRefresh time.Duration `mapstructure:"max-refresh" json:"max-refresh"`
}

func NewJwtOptions() *JwtOptions {
	return &JwtOptions{
		Realm:      "mxshop",
		Key:        "mxshop",
		Timeout:    24 * time.Hour,
		MaxRefresh: 24 * time.Hour,
	}
}

func (s *JwtOptions) Validate() []error {
	var errs []error

	//Key不能随便填，验证 Key 字段的长度是否在6到32个字符之间
	if !govalidator.StringLength(s.Key, "6", "32") {
		errs = append(errs, fmt.Errorf("--secret-key must larger than 5 and little than 33"))
	}

	return errs
}

func (s *JwtOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.StringVar(&s.Realm, "jwt.realm", s.Realm, "Realm name to display to the user.")
	fs.StringVar(&s.Key, "jwt.key", s.Key, "Private key used to sign jwt token.")
	fs.DurationVar(&s.Timeout, "jwt.timeout", s.Timeout, "JWT token timeout.")
	fs.DurationVar(&s.MaxRefresh, "jwt.max-refresh", s.MaxRefresh, ""+
		"This field allows clients to refresh their token until MaxRefresh has passed.")
}
