package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/patrickmn/go-cache"
	"github.com/sbrki/snote/internal/storage"
)

type Server struct {
	storage          storage.Storage
	templateRegistry *TemplateRegistry
	echo             *echo.Echo
	renderCache      *cache.Cache
}

func NewServer(storage storage.Storage, templateRegistry *TemplateRegistry) *Server {
	s := new(Server)
	s.storage = storage
	s.templateRegistry = templateRegistry
	s.echo = echo.New()
	s.echo.Renderer = s.templateRegistry
	s.renderCache = cache.New(1*time.Hour, 1*time.Minute)

	s.echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, " +
			"lat=${latency_human}, ip=${remote_ip}, in=${bytes_in}, out=${bytes_out}\n",
	}))

	s.echo.Static("/static", "web/static")
	s.echo.File("/favicon.ico", "web/static/favicon.ico")

	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	s.echo.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "edit.html", nil)
	})

	s.echo.GET("/:note_id/edit", func(c echo.Context) error {
		id := c.Param("note_id")
		_, err := s.storage.LoadNote(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "404 Not found")
		}

		return c.Render(http.StatusOK, "edit.html", nil)
	})

	s.echo.GET("/:note_id", func(c echo.Context) error {
		id := c.Param("note_id")
		note, err := s.storage.LoadNote(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "404 Not found")
		}

		// render the note markdown contents to html.
		// rendered html is cached as it takes approx. 1s to render a 5k LoC markdown.
		// first, check if html rendering of markdown exists in cache
		html, found := s.renderCache.Get(note.ID)
		if !found {
			// if not, render it
			parser := parser.NewWithExtensions(parser.CommonExtensions)
			html = markdown.ToHTML([]byte(note.Contents), parser, nil)
			// add it to cache
			s.renderCache.SetDefault(note.ID, html)
		}

		return c.Render(http.StatusOK, "preview.html", struct {
			RenderedHTML string
			ID           string
		}{fmt.Sprintf("%s", html), note.ID})
	})

	s.echo.GET("/somenote", func(c echo.Context) error {
		n := new(storage.Note)
		n.ID = "prvi"
		n.Contents = "# helo\n* lista\n* jos!\n	* podlista?"
		n.Title = "helo"
		n.IsPublic = false
		n.LastEdit = time.Now()

		//fmt.Println(n.ParseTitle())

		err := s.storage.SaveNote(n)
		if err != nil {
			c.Logger().Fatal(err)
		}

		note, err := s.storage.LoadNote("prvi")
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Println("Note read from disk")
		fmt.Println(note)

		return c.JSON(http.StatusOK, n)
	})

	// API
	s.echo.GET("/api/note/:note_id", s.noteGetHandler)
	s.echo.PUT("/api/note/:note_id", s.notePutHandler)
}

func (s *Server) Run() {
	s.echo.Logger.Fatal(s.echo.Start(":8081"))
}
