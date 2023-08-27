package websocket

import (
	"encoding/json"
	"fmt"
	"time"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChangeGame  = "change_game"

	GameOutcomes       = "gameprogress-event"
	LeaderboardChanges = "leaderboard-game-event"
)

type SendMessageEvent struct {
	Message   string `json:"message"`
	GameEvent string `json:"gameEvent"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeGameEvent struct {
	Name string `json:"name"`
}

func SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var sendMessageEvent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &sendMessageEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Prepare an Outgoing Message to others
	var broadMessage NewMessageEvent

	broadMessage.Sent = time.Now()
	broadMessage.Message = sendMessageEvent.Message
	broadMessage.GameEvent = sendMessageEvent.GameEvent
	c.game = sendMessageEvent.GameEvent

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage
	// Broadcast to all other Clients
	for client := range c.manager.clients {
		// Only send to clients subscribed to this event
		if client.game == c.game {
			client.egress <- outgoingEvent
		}

	}
	return nil
}

func ChangeGameHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var changeGameEvent ChangeGameEvent
	if err := json.Unmarshal(event.Payload, &changeGameEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// Add Client to game
	c.game = changeGameEvent.Name

	return nil
}
