package config

import "time"

type ObjectStorage struct {
	Endpoint          string        `mapstructure:"endpoint"`
	Bucket            string        `mapstructure:"bucket"`
	AccessKey         string        `mapstructure:"access_key"`
	SecretKey         string        `mapstructure:"secret_key"`
	UseSSL            bool          `mapstructure:"use_ssl"`
	PresignExpiration time.Duration `mapstructure:"presign_expiration"`
	MaxFileSize       int64         `mapstructure:"max_file_size"`
}
