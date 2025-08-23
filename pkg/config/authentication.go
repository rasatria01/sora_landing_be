package config

import "time"

type Authentication struct {
	EncryptKey         string        `mapstructure:"encrypt_key"`
	AccessSecretKey    string        `mapstructure:"access_secret_key"`
	RefreshSecretKey   string        `mapstructure:"refresh_secret_key"`
	Issuer             string        `mapstructure:"issuer"`
	AccessTokenExpiry  time.Duration `mapstructure:"access_token_expiry"`
	RefreshTokenExpiry time.Duration `mapstructure:"refresh_token_expiry"`
}
