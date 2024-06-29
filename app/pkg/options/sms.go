package options

import (
	"github.com/spf13/pflag"
)

type SmsOptions struct {
	APIKey    string `mapstructure:"APIKey" json:"APIKey"`
	APISecret string `mapstructure:"APISecret" json:"APISecret"`
}

func NewSmsOptions() *SmsOptions {
	return &SmsOptions{
		APIKey:    "",
		APISecret: "",
	}
}

func (s *SmsOptions) Validate() []error {
	var errs []error
	return errs
}

func (s *SmsOptions) AddFlags(fs *pflag.FlagSet) {
	if fs == nil {
		return
	}

	fs.StringVar(&s.APIKey, "sms.apikey", s.APIKey, "sms API key")
	fs.StringVar(&s.APISecret, "sms.apiSecret", s.APISecret, "sms API secret")
}
