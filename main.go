package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/gotailwindcss/tailwind/twembed"
	"github.com/gotailwindcss/tailwind/twhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(c echo.Context) error {

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()
	return nil
}

type team struct {
	Name   string
	id     int
	Points int
}

type teams struct {
	Teams []team
}

func createTeam(name string, id int) team {
	return team{
		Name:   name,
		id:     id,
		Points: 0,
	}
}

func NewTeams() *teams {
	return &teams{
		Teams: []team{},
	}
}

func main() {

	e := echo.New()
	e.Renderer = newTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/css", "/css")
	e.Static("/scripts", "/scripts")

	listofTeams := NewTeams()
	team_id := 0 //start the id teams

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.POST("/teams", func(c echo.Context) error {
		team_id++
		name := c.FormValue("name")
		team := createTeam(name, team_id)
		listofTeams.Teams = append(listofTeams.Teams, team)
		return c.Render(200, "team", team)
	})

	e.POST("/host/team/:id", func(c echo.Context) error {
		winnerteam := team{}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("not a number")
		}
		points := c.Param("points")
		pointstoappend, err := strconv.Atoi(points)
		if err != nil {
			fmt.Println("not a number")
		}
		for _, team := range listofTeams.Teams {
			if id == team.id {
				team.Points += pointstoappend
				winnerteam = team
			}
		}
		return c.Render(200, "points", winnerteam.Points)
	})

	e.POST("/host/team/:id", func(c echo.Context) error {
		winnerteam := team{}
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("not a number")
		}
		points := c.Param("points")
		pointstoappend, err := strconv.Atoi(points)
		if err != nil {
			fmt.Println("not a number")
		}
		for _, team := range listofTeams.Teams {
			if id == team.id {
				team.Points += pointstoappend
				winnerteam = team
			}
		}
		return c.Render(200, "points", winnerteam.Points)
	})

	tailwindHandler := twhandler.New(http.Dir("css"), "/css", twembed.New())
	e.GET("/css/*", echo.WrapHandler(tailwindHandler))
	e.GET("ws", handleWebSocket)

	port := ":8080"
	fmt.Println("WebSocket server is running on Port", port)
	e.Logger.Fatal(e.Start(port))
}
