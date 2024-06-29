package options

import (
	"github.com/spf13/pflag"
	"time"
)

type MySQLOptions struct {
	Host                  string        `mapstructure:"host" json:"host,omitempty"`
	Port                  int           `mapstructure:"port" json:"port,omitempty"`
	Username              string        `mapstructure:"username" json:"username,omitempty"`
	Password              string        `mapstructure:"password" json:"password,omitempty"`
	Database              string        `mapstructure:"database" json:"database,omitempty"`
	MaxIdleConnections    int           `mapstructure:"max-idle-connections" json:"max-idle-connections"`
	MaxOpenConnections    int           `mapstructure:"max-open-connections" json:"max-open-connections"`
	MaxConnectionLifetime time.Duration `mapstructure:"max-connection-life-time" json:"max-connection-life-time"`
	LogLevel              int           `mapstructure:"log-level" json:"log-level"`
}

func NewMySQLOptions() *MySQLOptions {
	return &MySQLOptions{
		Host:                  "127.0.0.1",
		Port:                  3306,
		Username:              "",
		Password:              "",
		Database:              "",
		MaxIdleConnections:    10,
		MaxOpenConnections:    100,
		MaxConnectionLifetime: time.Duration(10) * time.Second,
		LogLevel:              1, //silent
	}
}

// Validate verifies flags passed to ServerOptions.
func (so *MySQLOptions) Validate() []error {
	errs := []error{}
	return errs
}

// AddFlags adds flags related to server storage for a specific APIServer to the specified FlagSet.
func (mo *MySQLOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&mo.Host, "mysql.host", mo.Host, ""+
		"MySQL service host address. If left blank, the following related mysql options will be ignored.")

	fs.IntVar(&mo.Port, "mysql.port", mo.Port, ""+
		"MySQL service port")

	fs.StringVar(&mo.Username, "mysql.username", mo.Username, ""+
		"Username for access to mysql service.")

	fs.StringVar(&mo.Password, "mysql.password", mo.Password, ""+
		"Password for access to mysql, should be used pair with password.")

	fs.StringVar(&mo.Database, "mysql.database", mo.Database, ""+
		"Database name for the server to use.")

	fs.IntVar(&mo.MaxIdleConnections, "mysql.max-idle-connections", mo.MaxOpenConnections, ""+
		"Maximum idle connections allowed to connect to mysql.")

	fs.IntVar(&mo.MaxOpenConnections, "mysql.max-open-connections", mo.MaxOpenConnections, ""+
		"Maximum open connections allowed to connect to mysql.")

	fs.DurationVar(&mo.MaxConnectionLifetime, "mysql.max-connection-life-time", mo.MaxConnectionLifetime, ""+
		"Maximum connection life time allowed to connecto to mysql.")

	fs.IntVar(&mo.LogLevel, "mysql.log-mode", mo.LogLevel, ""+
		"Specify gorm log level.")
}
