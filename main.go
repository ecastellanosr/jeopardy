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

func main() {

	e := echo.New()
	e.Renderer = newTemplate()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/css", "/css")
	e.Static("/scripts", "/scripts")

	listofTeams := NewTeams()
	team_id := 0 //start the id teams

	categories, err := readcategories()
	if err != nil {
		fmt.Println(err)
	}

	e.GET("/", func(c echo.Context) error {
		team_id = 0 // take this out before solution in websockets
		return c.Render(200, "index", categories)
	})

	e.GET("/host", func(c echo.Context) error {
		return c.Render(200, "host", categories)
	})

	e.POST("/question", func(c echo.Context) error {
		number := c.FormValue("number")
		fmt.Println(number)
		return c.Render(200, "clicked-card", number)
	})

	e.POST("/teams", func(c echo.Context) error {
		team_id++

		name := c.FormValue("name")
		team := createTeam(name, team_id)
		listofTeams.Teams = append(listofTeams.Teams, team)

		c.Render(200, "team", team)
		if team_id >= 4 {
			// take out the question shield
			return c.Render(200, "oob-cover", nil)
		}
		return c.Render(200, "oob-add-team", nil)
	})

	e.POST("/yes-team", func(c echo.Context) error {
		return c.Render(200, "team-form", nil)

	})

	e.POST("/no-team", func(c echo.Context) error {
		if team_id > 1 {
			// take out the question shield
			return c.Render(200, "oob-cover", nil)
		}
		return c.Render(200, "add-team", nil)
	})

	e.POST("/host/team/:team_id/:question_id", func(c echo.Context) error {
		team_id := c.Param("team_id")
		team_id_int, err := strconv.Atoi(team_id)
		if err != nil {
			fmt.Println("team_id is not a number")
		}
		question_id := c.Param("team_id")
		question_id_int, err := strconv.Atoi(question_id)
		if err != nil {
			fmt.Println("question_id is not a number")
		} // case

		winnerteam := addpointstoteam(categories, listofTeams, question_id_int, team_id_int)
		return c.Render(200, "oob-points", winnerteam)
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
			if id == team.ID {
				team.Points += pointstoappend
				winnerteam = team
			}
		}
		return c.Render(200, "points", winnerteam.Points)
	})

	tailwindHandler := twhandler.New(http.Dir("css"), "/css", twembed.New())
	e.GET("/css/*", echo.WrapHandler(tailwindHandler))
	e.GET("ws", handleWebSocket)

	port := ":8000"
	fmt.Println("WebSocket server is running on Port", port)
	e.Logger.Fatal(e.Start(port))
}
