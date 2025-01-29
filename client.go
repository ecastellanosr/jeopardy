package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second
	//send pings to peer with this time period. less than pong
	pingPeriod = ((pongWait * 9) / 10)
	//time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

type Message struct {
	Text   string
	Type   string
	Client Client
}

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager    *Manager
	host       bool
	egress     chan []byte
	id         string
	Status     string
}

type Payload struct {
	Name    string `json:"name"`
	HEADERS struct {
		HXRequest     string      `json:"HX-Request"`
		HXTrigger     interface{} `json:"HX-Trigger"`
		HXTriggerName string      `json:"HX-Trigger-Name"`
		HXTarget      string      `json:"HX-Target"`
		HXCurrentURL  string      `json:"HX-Current-URL"`
	} `json:"HEADERS"`
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	NewID := uuid.New()
	NewIDstring := NewID.String()
	return &Client{
		connection: conn,
		manager:    manager,
		host:       false,
		egress:     make(chan []byte),
		id:         NewIDstring,
		Status:     "client",
	}
}

func NewHost(conn *websocket.Conn, manager *Manager) *Client {
	NewID := uuid.New()
	NewIDstring := NewID.String()
	return &Client{
		connection: conn,
		manager:    manager,
		host:       true,
		egress:     make(chan []byte),
		id:         NewIDstring,
		Status:     "host",
	}
}

func (c *Client) readMessages() {
	defer func() {
		//cleanup connection
		c.manager.removeClient(c)
	}()
	c.connection.SetReadLimit(maxMessageSize)
	c.connection.SetReadDeadline(time.Now().Add(pongWait))
	c.connection.SetPongHandler(func(appData string) error {
		c.connection.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseAbnormalClosure, websocket.CloseGoingAway) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		log.Println("messageType", messageType)
		log.Println("payload:", string(payload))
		var mPayload Payload
		if err = json.Unmarshal(payload, &mPayload); err != nil {
			log.Printf("error marshaling message: %v", err)
		}

		log.Println("payload name:", mPayload.Name)
		// send the message to the broadcast channel.
		c.manager.broadcast <- &Message{Text: mPayload.Name, Type: mPayload.HEADERS.HXTarget, Client: *c}

	}
}

func (c *Client) writeMessages() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//cleanup connection
		c.manager.removeClient(c)
	}()
	for {
		select {
		case message, ok := <-c.egress:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				err := c.connection.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Println("connection closed: ", err)
				}
				return
			}
			w, err := c.connection.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("failed to send the message: %v", err)
			}
			w.Write(message)

			n := len(c.egress)
			for i := 0; i < n; i++ {
				w.Write(message)
			}

			if err := w.Close(); err != nil {
				log.Printf("problem while closing")
				return
			}
		case <-ticker.C:
			c.connection.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
