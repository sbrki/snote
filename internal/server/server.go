package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sbrki/snote/internal/storage"
)

type Server struct {
	storage          storage.Storage
	templateRegistry *TemplateRegistry
	echo             *echo.Echo
}

func NewServer(storage storage.Storage, templateRegistry *TemplateRegistry) *Server {
	s := new(Server)
	s.storage = storage
	s.templateRegistry = templateRegistry
	s.echo = echo.New()
	s.echo.Renderer = s.templateRegistry

	s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, " +
			"lat=${latency_human}, ip=${remote_ip}, in=${bytes_in}, out=${bytes_out}\n",
	}))

	return s
}

func (s *Server) Run() {
	// web stuff
	e := echo.New()

	e.Static("/static", "web/static")
	e.File("/favicon.ico", "web/static/favicon.ico")

	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "edit.html", nil)
	})

	e.GET("/user/:name", func(c echo.Context) error {
		name := c.Param("name")
		return c.String(http.StatusOK, name)
	})

	e.GET("/somenote", func(c echo.Context) error {
		n := new(storage.Note)
		n.ID = "prvi"
		n.Contents = "#helo"
		n.Title = "helo"
		n.IsPublic = false
		n.LastEdit = time.Now()

		//fmt.Println(n.ParseTitle())

		err := s.storage.SaveNote(n)
		if err != nil {
			c.Logger().Fatal(err)
		}

		err, note := s.storage.LoadNote("prvi")
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Println("Note read from disk")
		fmt.Println(note)

		return c.JSON(http.StatusOK, n)
	})

	e.Logger.Fatal(e.Start(":8081"))

}
