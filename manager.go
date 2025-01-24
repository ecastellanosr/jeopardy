package main

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Manager struct {
	clients   ClientList
	broadcast chan *Message
	teams     []team
	messages  []*Message
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients:   make(ClientList),
		broadcast: make(chan *Message),
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (m *Manager) handleWebSocket(c echo.Context) error {

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	client := NewClient(conn, m)
	m.addClient(client)
	log.Println("new client")
	go client.readMessages()
	go client.writeMessages()
	go client.manager.ListenBroadcast(c)
	return nil
}

func (m *Manager) handleHostWebSocket(c echo.Context) error {

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	client := NewHost(conn, m)
	m.addClient(client)
	log.Println("new host")
	go client.readMessages()
	go client.writeMessages()
	go client.manager.ListenBroadcast(c)
	return nil
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}

func (m *Manager) ListenBroadcast(c echo.Context) {
	team_id := 0
	for {
		select {
		case msg := <-m.broadcast:
			if msg.Type == "team-form" {
				team := createTeam(msg.Text, team_id)
				m.teams = append(m.teams, team)
				for client := range m.clients {
					if client.host == true {
						select {
						case client.egress <- addTeamTemplate(&m.teams[team_id]):
							team_id++
						default:
							close(client.egress)
							m.removeClient(client)
						}
					}
				}
			}

		}
	}
}

func addTeamTemplate(team *team) []byte {
	tmpl, err := template.ParseFiles("views/team.html")
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	var renderedMessage bytes.Buffer
	err = tmpl.Execute(&renderedMessage, team)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	log.Printf("this got to here, before sending the bytes to the host")
	return renderedMessage.Bytes()
}
