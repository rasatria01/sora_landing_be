package config

type Logger struct {
	Environment string `mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level"`
	Encoding    string `mapstructure:"encoding"`
}
