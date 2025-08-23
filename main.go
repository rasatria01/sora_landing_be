package main

import (
	"log"
	"sora_landing_be/cmd/routes"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/config"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/http/server"
	"sora_landing_be/pkg/logger"
)

func main() {
	cfg := config.LoadConfig()

	logger.NewZapLogger(cfg.Logger)

	database.InitDB(cfg.Database)

	authentication.NewJWTManager(authentication.JWTOptions{
		AccessSecret:       cfg.Authentication.AccessSecretKey,
		RefreshSecret:      cfg.Authentication.RefreshSecretKey,
		Issuer:             cfg.Authentication.Issuer,
		ExpiryAccessToken:  cfg.Authentication.AccessTokenExpiry,
		ExpiryRefreshToken: cfg.Authentication.RefreshTokenExpiry,
	})

	authentication.SetupKey(cfg.Authentication.EncryptKey)

	server.Init(cfg.Application, routes.RegisterV1).GracefulShutdown()
	defer func() {
		err := logger.Log.Sync()
		if err != nil {
			log.Println(err)
		}
	}()
}
