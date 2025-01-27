package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Manager struct {
	clients            ClientList
	broadcast          chan *Message
	categoriesandteams categoriesandteams
	teams              []team
	categories         []Category
	messages           []*Message
	currentcard        int
	currentTeam        team
	currentTeamID      int
	answerRevealed     bool
	sync.RWMutex
}

func NewManager(c []Category) *Manager {
	return &Manager{
		clients:   make(ClientList),
		broadcast: make(chan *Message),
		categoriesandteams: categoriesandteams{
			Categories: c,
			Fullteams:  false,
		},
		categories:     c,
		currentcard:    0,
		currentTeamID:  1,
		answerRevealed: false,
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
	// for i, team := range m.teams {
	// 	log.Println("adding team:", team.Name)
	// 	client.egress <- addTeamTemplate(&m.teams[i])
	// }
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
	// for i, team := range m.teams {
	// 	log.Println("adding team:", team.Name)
	// 	client.egress <- addTeamTemplate(&m.teams[i])
	// }
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
		log.Println("a user was deleted")
		delete(m.clients, client)
	}
}

func (m *Manager) ListenBroadcast(c echo.Context) {
	for {
		select {
		case msg := <-m.broadcast:
			if msg.Type == "team-form" {
				if m.currentTeamID >= 5 {
					continue
				}
				team := createTeam(msg.Text, m.currentTeamID)
				m.teams = append(m.teams, team)
				m.categoriesandteams.Teams = append(m.categoriesandteams.Teams, &team)
				m.currentTeamID++
				if m.currentTeamID == 4 {
					m.categoriesandteams.Fullteams = true
				}
				for client := range m.clients {
					if client.host != true {
						continue
					}
					select {
					case client.egress <- addTeamTemplate(&m.teams[(team.ID - 1)]):
						log.Println("this is the name of the team that was sent:", m.teams[(team.ID-1)].Name)
					default:
						close(client.egress)
						m.removeClient(client)
					}

				}

			}
			if strings.HasPrefix(msg.Type, "card") {
				cardN, _ := strings.CutPrefix(msg.Type, "card")
				id, err := strconv.Atoi(cardN)
				if err != nil {
					log.Printf("error while converting the card number in wsmessage, %s", err)
					continue
				}
				card, category_id := FindCard(m.categories, id)
				for i, revealedcard := range m.categories[category_id].Cards {
					if card.ID == revealedcard.ID {
						m.categories[category_id].Cards[i].Revealed = true
					}
				}
				for client := range m.clients {
					if client.host == true {
						HostCard := CardSelection{
							ClientStatus: "host",
							ID:           card.ID,
							Number:       card.Number,
						}
						select {
						case client.egress <- showSelectedCard(HostCard):
						default:
							close(client.egress)
							m.removeClient(client)
						}
					}
					if client.host == false {
						ClientCard := CardSelection{
							ClientStatus: "client",
							ID:           card.ID,
							Number:       card.Number,
						}
						m.currentcard = card.ID //as the new card is being selected the current card is changed
						select {
						case client.egress <- showSelectedCard(ClientCard):
						default:
							close(client.egress)
							m.removeClient(client)
						}
					}
				}
			}
			if msg.Type == "reveal-button" {
				if m.currentcard == 0 {
					continue
				}
				m.answerRevealed = true
				card, _ := FindCard(m.categories, m.currentcard)
				for client := range m.clients {
					if client.host == true {
						continue
					}
					select {
					case client.egress <- showCardAnswer(card):
					default:
						close(client.egress)
						m.removeClient(client)
					}
				}
			}
			if strings.HasPrefix(msg.Type, "addpoints-") {
				if m.currentcard == 0 || m.answerRevealed == false {
					continue
				}
				teamID_string, _ := strings.CutPrefix(msg.Type, "addpoints-")
				teamID, err := strconv.Atoi(teamID_string)
				if err != nil {
					log.Println("this is the teamidstring cut:", teamID_string)
					log.Printf("error while converting the team ID to an int in wsmessage, %s", err)
					continue
				}
				for i, team := range m.teams {
					if team.ID != teamID {
						continue
					}
					card, _ := FindCard(m.categories, m.currentcard)
					cardpoints, err := strconv.Atoi(card.Number)
					if err != nil {
						log.Printf("error while converting the card number points to an int in wsmessage, %s", err)
						continue
					}
					m.teams[i].Points += cardpoints
					m.categoriesandteams.Teams[i].Points = m.teams[i].Points
					log.Println("current points of team:", m.teams[i].Points, "cardpoints:", cardpoints)
					m.currentTeam = m.teams[i]
					log.Println("current points of currentteam:", m.currentTeam.Points, m.currentTeam.Name, m.currentTeamID)
				}
				for client := range m.clients {
					select {
					case client.egress <- addpoints(m.currentTeam):
						client.egress <- resetQanimation()
						client.egress <- removeQuestionCover()
					// case client.egress <- resetQanimation():
					default:
						close(client.egress)
						m.removeClient(client)
					}
				}
				m.currentTeam = team{
					Name:   "",
					ID:     0,
					Points: 0,
				}
				m.currentcard = 0
				m.answerRevealed = false
			}
		}
	}
}
