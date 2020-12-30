package main

import (
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/labstack/echo"
	"github.com/sbrki/snote/internal/server"
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
	diskStoragePath, isSet := os.LookupEnv("STORAGE_PATH")
	if !isSet {
		diskStoragePath = "/tmp"
	}
	fmt.Println("Using storage path:", diskStoragePath)
	ds := storage.NewDiskStorage(diskStoragePath)
	// setup templates
	tr := server.NewTemplateRegistry("web/templates/*.html")
	// create server
	serv := server.NewServer(ds, tr)
	serv.Run()
}
