package server

import (
	"fmt"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo"
	"github.com/sbrki/snote/internal/storage"
)

func (s *Server) htmlIndexHandler(c echo.Context) error {
	return c.Redirect(http.StatusTemporaryRedirect, "/ls")
}

func (s *Server) htmlNoteHandler(c echo.Context) error {
	id := c.Param("note_id")
	var note *storage.Note

	if id == "ls" || id == "lstag" {
		note := new(storage.Note)
		if id == "ls" {
			note.GenerateLs(s.storage)
		} else {
			note.GenerateLsTag(s.storage)
		}
		parser := parser.NewWithExtensions(parser.CommonExtensions)
		html := markdown.ToHTML([]byte(note.Contents), parser, nil)
		return c.Render(http.StatusOK, "preview.html", struct {
			RenderedHTML string
			ID           string
		}{fmt.Sprintf("%s", html), note.ID})

	} else {
		storedNote, err := s.storage.LoadNote(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "404 Not found")
		}
		note = storedNote
	}

	fmt.Println("TITLE:", note.ParseTitle())

	// render the note markdown contents to html.
	// rendered html is cached as it takes approx. 1s to render a 5k LoC markdown.
	// first, check if html rendering of markdown exists in cache
	html, found := s.renderCache.Get(note.ID)
	if !found {
		// if not, render it
		html = note.RenderHTML()
		// add it to cache
		s.renderCache.SetDefault(note.ID, html)
	}

	return c.Render(http.StatusOK, "preview.html", struct {
		RenderedHTML string
		ID           string
	}{fmt.Sprintf("%s", html), note.ID})
}

func (s *Server) htmlNoteEditHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "edit.html", nil)
}
