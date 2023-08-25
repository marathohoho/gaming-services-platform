package main

import (
	"context"
	"gaming-services-platform/config"
	"gaming-services-platform/grpc"
	"gaming-services-platform/internal/repositories"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg := config.Init()

	rdb := redis.NewClient(&redis.Options{
		Addr:       cfg.Redis.Addr,
		DB:         cfg.Redis.Db,
		MaxRetries: 3,
	})

	walletRepository := repositories.NewWalletRepository(context.Background(), rdb)
	grpcService := grpc.NewGrpcService(rdb, walletRepository, cfg.GrpcServer)

	grpcService.ListenForConnection(context.Background())
}
