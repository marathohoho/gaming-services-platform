package main

import (
	"context"
	"gaming-services-platform/config"
	"gaming-services-platform/internal/repositories"
	"gaming-services-platform/server"
	"gaming-services-platform/user"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg := config.Init()
	userServer := server.Init()

	rdb := redis.NewClient(&redis.Options{
		Addr:       cfg.Redis.Addr,
		DB:         cfg.Redis.Db,
		MaxRetries: 3,
	})

	userRepository := repositories.NewUserRepository(context.Background(), rdb)
	userService := user.NewUserService(rdb, userRepository)

	userServer.Post("/users", userService.Create())
	userServer.Get("/user/:id", userService.Get())

	go server.Listen(userServer, cfg.UserServer)
	server.Shutdown(userServer, nil)
}
