package main

import (
	"fmt"
	"io"
	"net/http"
	"text/template"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/sbrki/snote/internal/storage"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

var t = &Template{
	templates: template.Must(template.ParseGlob("web/templates/*.html")),
}

func main() {
	// setup disk storage
	ds := storage.NewDiskStorage("/tmp/snotestorage")

	// web stuff
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}, " +
			"lat=${latency_human}, ip=${remote_ip}, in=${bytes_in}, out=${bytes_out}\n",
	}))
	e.Renderer = t // register templates

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

		err := ds.SaveNote(n)
		if err != nil {
			c.Logger().Fatal(err)
		}

		err, note := ds.LoadNote("prvi")
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Println("Note read from disk")
		fmt.Println(note)

		return c.JSON(http.StatusOK, n)
	})

	e.Logger.Fatal(e.Start(":8081"))
}
