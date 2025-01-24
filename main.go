package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gotailwindcss/tailwind/twembed"
	"github.com/gotailwindcss/tailwind/twhandler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// html template struct to send html instead of json
type Template struct {
	tmpl *template.Template
}

// initialize a new template
func newTemplate() *Template {
	return &Template{
		tmpl: template.Must(template.ParseGlob("views/*.html")),
	}
}

// Render a new template to the html that made the request, usually just a short part of the html
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.tmpl.ExecuteTemplate(w, name, data)
}

// Render with a websocket call instead of an ajax one
func (t *Template) WSRender(name string, data interface{}, c echo.Context) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

func main() {
	// initialize echo
	e := echo.New()
	e.Renderer = newTemplate()

	// initialize a user manager
	manager := NewManager()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// let server use the css, scripts and imgs folders
	e.Static("/css", "./css")
	e.Static("/scripts", "./scripts")
	e.Static("/imgs", "./imgs")
	// initialize a new list of teams
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

	e.POST("/question/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Render(500, "nothing", nil)
		}
		card := FindCard(categories, id)
		return c.Render(200, "deletedquestion", card)
	})

	e.POST("/revealquestion/:id", func(c echo.Context) error {
		// time.Sleep(1 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Render(500, "nothing", nil)
		}
		card := FindCard(categories, id)
		return c.Render(200, "question-cover", card)
	})

	e.POST("/teams", func(c echo.Context) error {
		team_id++

		name := c.FormValue("name")
		team := createTeam(name, team_id)
		listofTeams.Teams = append(listofTeams.Teams, team)

		// // function that sends a message to all host == true clients
		// for client := range manager.clients {
		// 	if client.host == true {
		// 		log.Println("posting teams to host")
		// 		tmpl := c.Echo().Renderer.(*Template)
		// 		rendered, err := tmpl.WSRender("test-host-team", team, c)
		// 		if err != nil {
		// 			log.Printf("could not send the team name to the host")
		// 		}
		// 		err = client.connection.WriteMessage(websocket.TextMessage, rendered)
		// 		if err != nil {
		// 			log.Printf("could not send the team name to the host")
		// 		}

		// 	}
		// }
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

	e.POST("/host/addpoints/:team_id/:question_id", func(c echo.Context) error {
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
	e.GET("ws", manager.handleWebSocket)
	e.GET("wshost", manager.handleHostWebSocket)

	port := ":8000"
	fmt.Println("WebSocket server is running on Port", port)
	e.Logger.Fatal(e.Start(port))
}
