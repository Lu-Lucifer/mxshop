package options

import "github.com/spf13/pflag"

type NacosOptions struct {
	Host        string `mapstructure:"host" json:"host"`
	Port        uint64 `mapstructure:"port" json:"port"`
	NamespaceId string `mapstructure:"namespaceId" json:"namespaceId"`
	DataId      string `mapstructure:"dataId" json:"dataId"`
	Group       string `mapstructure:"group" json:"group"`
	UserName    string `mapstructure:"userName" json:"userName"`
	Password    string `mapstructure:"password" json:"password"`
}

func NewNacosOptions() *NacosOptions {
	return &NacosOptions{
		Host:        "127.0.0.1",
		Port:        8848,
		NamespaceId: "public",
		DataId:      "flow",
		Group:       "sentinel-go",
		UserName:    "nacos",
		Password:    "nacos",
	}
}

func (n *NacosOptions) Validate() []error {
	errs := []error{}
	return errs
}

// AddFlags adds flags related to server storage for a specific APIServer to the specified FlagSet.
func (n *NacosOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&n.Host, "nacos.host", n.Host, "nacos host")

	fs.Uint64Var(&n.Port, "nacos.port", n.Port, "nacos port")

	fs.StringVar(&n.NamespaceId, "nacos.NamespaceId", n.NamespaceId, "nacos NamespaceId")

	fs.StringVar(&n.DataId, "nacos.dataId", n.DataId, "nacos dataid")

	fs.StringVar(&n.Group, "nacos.group", n.Group, "nacos group")
	fs.StringVar(&n.UserName, "nacos.UserName", n.Group, "nacos UserName")
	fs.StringVar(&n.Password, "nacos.Password", n.Group, "nacos Password")
}
