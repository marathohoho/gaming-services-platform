package main

import (
	"gaming-services-platform/websocket"
	"log"
	"net/http"
)

func main() {
	setupAPI()

	log.Fatal(http.ListenAndServe(":5104", nil))
}

func setupAPI() {
	manager := websocket.NewManager()

	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServeWS)
}
