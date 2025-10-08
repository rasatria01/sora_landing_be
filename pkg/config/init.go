package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

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
	var err error
	once.Do(func() {
		// Use CONFIG_PATH env if set (Docker friendly), else default path
		configPath := os.Getenv("CONFIG_PATH")
		if configPath == "" {
			configPath = "./pkg/config/files"
		}

		viper.SetConfigName("env") // name of your config file without extension
		viper.SetConfigType("yaml")
		viper.AddConfigPath(configPath)

		// Set up environment variable bindings
		viper.SetEnvPrefix("")
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()

		// Read config file
		if err = viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				err = fmt.Errorf("config file not found in %s", configPath)
				return
			}
			err = fmt.Errorf("error reading config file: %w", err)
			return
		}

		// Check for PORT environment variable specifically
		if port := os.Getenv("PORT"); port != "" {
			portNum, err := strconv.Atoi(port)
			if err == nil {
				viper.Set("application.port", portNum)
			}
		}

		// Unmarshal into struct
		if err = viper.Unmarshal(&config); err != nil {
			err = fmt.Errorf("error unmarshaling config: %w", err)
			return
		}
	})

	if err != nil {
		log.Fatalf("error loading config: %s", err.Error())
	}

	return config
}
