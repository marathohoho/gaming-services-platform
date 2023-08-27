package websocket

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	egress     chan Event // to avoid concurrent writes on ws
	game       string
}

type ClientList map[*Client]bool

func (c *Client) readMessages() {
	defer func() {
		// close connection once done reading messages
		c.manager.removeClient(c)
	}()

	for {
		_, payload, err := c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unable to read message: %v", err)
			}

			break
		}

		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("unable to unmarshal message: %v", err)
		}

		if err := c.manager.RouteEvent(request, c); err != nil {
			log.Println("unable to handle the event: ", err)
		}

	}
}

func (c *Client) writeMessages() {
	defer func() {
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			if ok == false {
				if err := c.connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("connection closed: ", err)
				}

				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return
			}
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}

			log.Println("message was sent")
		}
	}
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
		egress:     make(chan Event),
		game:       GameOutcomes,
	}
}
