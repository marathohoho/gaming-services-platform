package main

import (
	"gaming-services-platform/config"
	"gaming-services-platform/server"
	"gaming-services-platform/user"
	"gaming-services-platform/wallet"
)

func main() {
	cfg := config.Init()
	app := server.Init()
	api := app.Group("/api")

	userHandler := user.NewUserHandler(cfg.UserServerHost)
	walletHandler := wallet.NewWalletHandler(cfg.WalletServerHost)

	api.Post("/users", userHandler.Register())
	api.Get("/user/:id", userHandler.Get())

	walletGroup := api.Group("/wallet")
	walletGroup.Post("/deposit", walletHandler.Deposit())
	walletGroup.Post("/withdraw", walletHandler.Withdraw())

	// another endpoint for handling gRPC calls

	go server.Listen(app, cfg.ApiServerHost)
	server.Shutdown(app, nil)

}
