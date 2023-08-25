package main

import (
	"context"
	"gaming-services-platform/config"
	"gaming-services-platform/internal/repositories"
	"gaming-services-platform/server"
	"gaming-services-platform/wallet"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg := config.Init()
	walletServer := server.Init()

	rdb := redis.NewClient(&redis.Options{
		Addr:       cfg.Redis.Addr,
		DB:         cfg.Redis.Db,
		MaxRetries: 3,
	})

	walletRepository := repositories.NewWalletRepository(context.Background(), rdb)
	walletService := wallet.NewWalletService(rdb, walletRepository)

	walletServer.Post("/deposit", walletService.DepositFunds())
	walletServer.Post("/withdraw", walletService.WithdrawFunds())

	// this will be a gRPC
	// walletServer.Get()

	go server.Listen(walletServer, cfg.WalletServer)
	server.Shutdown(walletServer, nil)
}
