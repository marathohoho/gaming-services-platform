package main

import (
	"context"
	"gaming-services-platform/config"
	"gaming-services-platform/internal/repositories"
	"gaming-services-platform/server"
	"gaming-services-platform/wallet"
	"log"
	"net/url"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

func main() {
	cfg := config.Init()
	walletServer := server.Init()
	log.Print("started wallet service")
	rdb := redis.NewClient(&redis.Options{
		Addr:       cfg.Redis.Addr,
		DB:         cfg.Redis.Db,
		MaxRetries: 3,
	})

	websocketURL := "ws://" + cfg.WebsocketServer + "/ws"
	u, err := url.Parse(websocketURL)
	if err != nil {
		log.Fatal("Unable to parse websocket URL: ", err)
	}
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Unable to connect to websocket server: ", err)
	}
	defer connection.Close()

	walletRepository := repositories.NewWalletRepository(context.Background(), rdb)
	walletService := wallet.NewWalletService(rdb, walletRepository, connection)

	walletServer.Post("/deposit", walletService.DepositFunds())
	walletServer.Post("/withdraw", walletService.WithdrawFunds())

	go server.Listen(walletServer, cfg.WalletServer)
	server.Shutdown(walletServer, nil)
}
