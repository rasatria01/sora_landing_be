package config

type ApplicationEnvironment string

const (
	Development ApplicationEnvironment = "development"
	Production  ApplicationEnvironment = "production"
	Test        ApplicationEnvironment = "test"
)
