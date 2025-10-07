package main

import (
	"sora_landing_be/cmd/routes"
	"sora_landing_be/pkg/authentication"
	"sora_landing_be/pkg/config"
	"sora_landing_be/pkg/database"
	"sora_landing_be/pkg/http/server"
	"sora_landing_be/pkg/logger"

	"go.uber.org/zap"
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

	// Initialize the server
	srv := server.Init(cfg.Application, routes.RegisterV1)

	// Start listening for shutdown signals in a separate goroutine
	go func() {
		srv.GracefulShutdown()
	}()

	// Log that we're starting
	logger.Log.Info("Server is running", zap.Int("port", cfg.Application.Port))

	// Keep the main goroutine running
	select {}
}
