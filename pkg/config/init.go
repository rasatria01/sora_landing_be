package config

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

const configFilePath = "./pkg/config/files"

type Config struct {
	Application    Application    `yaml:"application"`
	Authentication Authentication `yaml:"authentication"`
	Database       Database       `yaml:"database"`
	Logger         Logger         `yaml:"logger"`
	ObjectStorage  ObjectStorage  `yaml:"object_storage"`
}

var once sync.Once
var config Config

func LoadConfig() Config {
	var (
		err error
	)
	once.Do(func() {
		// Load the configuration from the file
		// This is a placeholder for actual loading logic
		viper.SetConfigName("env")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configFilePath)

		viper.AutomaticEnv()
		if err = viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				err = fmt.Errorf("config file not found")
				return
				// Handle the case where the config file is not found
				// For example, you might want to log this or set default values
			}
			err = fmt.Errorf("error reading config file: %w", err)
			return
		}
		if err = viper.Unmarshal(&config); err != nil {
			return
		}
	})

	if err != nil {
		log.Fatalf(fmt.Sprintf("error loading config: %s", err.Error()))
	}
	return config
}
