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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// let server use the css, scripts and imgs folders
	e.Static("/css", "./css")
	e.Static("/scripts", "./scripts")
	e.Static("/imgs", "./imgs")
	// initialize a new list of teams
	team_id := 0 //start the id teams

	categories, err := readcategories()
	// initialize a manager for Users and hosts
	manager := NewManager(categories)
	if err != nil {
		fmt.Println(err)
	}
	// this is the initial get for normal clients, it returns the index block and is filled with the categories and teams that were previously set.
	// categories are filled while initializing a manager and teams are set in the client UI. If you refresh you get those teams.
	// a button to reset the teams may be adviced.
	e.GET("/", func(c echo.Context) error {
		team_id = 0 // take this out before solution in websockets
		return c.Render(200, "index", manager.categoriesandteams)
	})
	// initial get for hosts.
	e.GET("/host", func(c echo.Context) error {
		return c.Render(200, "host", manager.categoriesandteams)
	})
	// Post to get the teams request for clients. Host are handled in websockets
	e.POST("/teams", func(c echo.Context) error {
		//this can be resolved by websockets instead as it might create race condition with the team id created in the manager side.
		team_id++
		name := c.FormValue("name")
		team := createTeam(name, team_id)
		c.Render(200, "team", team)
		if manager.currentTeamID > 4 {
			// take out the question shield
			return c.Render(200, "oob-cover", nil)
		}
		return c.Render(200, "oob-add-team", nil)
	})
	//Post for clicking yes when creating a team in the client UI
	e.POST("/yes-team", func(c echo.Context) error {
		return c.Render(200, "team-form", nil)
	})
	// Post for clicking no when creating a team in the client UI
	e.POST("/no-team", func(c echo.Context) error {
		if manager.currentTeamID > 2 {
			//this boolean makes so that there is no more "do you want to add another team" element when refreshed
			manager.categoriesandteams.Fullteams = true
			// take out the question shield
			return c.Render(200, "oob-cover", nil)
		}
		return c.Render(200, "add-team", nil)
	})
	//
	e.POST("/clientquestion/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Render(500, "nothing", nil)
		}
		card, _ := FindCard(categories, id)
		return c.Render(200, "deletedquestion", card)
	})

	e.POST("/clientrevealquestion/:id", func(c echo.Context) error {
		// time.Sleep(1 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Render(500, "nothing", nil)
		}
		card, _ := FindCard(categories, id)
		return c.Render(200, "question-cover", card)
	})

	e.POST("/hostquestion/:id", func(c echo.Context) error {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Render(500, "nothing", nil)
		}
		card, _ := FindCard(categories, id)
		return c.Render(200, "hostdeletedquestion", card)
	})

	e.POST("/hostrevealquestion/:id", func(c echo.Context) error {
		// time.Sleep(1 * time.Second)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Render(500, "nothing", nil)
		}
		card, _ := FindCard(categories, id)
		return c.Render(200, "host-question-cover", card)
	})
	// Delete the question cover triggered when ending animation to reset the animation position
	e.DELETE("/QuestionCover", func(c echo.Context) error {
		// time.Sleep(1 * time.Second)
		return c.NoContent(http.StatusOK)
	})
	//tailwind
	tailwindHandler := twhandler.New(http.Dir("css"), "/css", twembed.New())
	e.GET("/css/*", echo.WrapHandler(tailwindHandler))
	// manage websockets for clients
	e.GET("ws", manager.handleWebSocket)
	// manage websockets for hosts
	e.GET("wshost", manager.handleHostWebSocket)
	//Used port
	port := ":8000"
	fmt.Println("WebSocket server is running on Port", port)
	e.Logger.Fatal(e.Start(port))
}
