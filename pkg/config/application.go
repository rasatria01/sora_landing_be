package config

type Application struct {
	Port        int                    `yaml:"port"`
	Environment ApplicationEnvironment `yaml:"environment"`
}
