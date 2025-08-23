package config

import "time"

type Database struct {
	Host            string        `mapstructure:"host"`
	Port            string        `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Name            string        `mapstructure:"name"`
	SSLMode         string        `mapstructure:"ssl_mode"`
	MaxOpenIdleConn int           `mapstructure:"max_open_idle_conn"`
	MaxOpenConn     int           `mapstructure:"max_open_conn"`
	MaxIdleConn     time.Duration `mapstructure:"max_idle_conn"`
}
