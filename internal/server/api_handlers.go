package server

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/sbrki/snote/internal/storage"
)

func (s *Server) noteGetHandler(c echo.Context) error {
	id := c.Param("note_id")
	note, err := s.storage.LoadNote(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "404 Not found")
	}
	return c.JSON(http.StatusOK, note)
}

func (s *Server) notePutHandler(c echo.Context) error {
	id := c.Param("note_id")
	_, err := s.storage.LoadNote(id)
	if err != nil {
		// TODO(sbrki): if note does not exist, create one on a PUT request
		return echo.NewHTTPError(http.StatusNotFound, "404 Not found")
	}
	// note exists - update its contents
	updatedNote := new(storage.Note)
	if err := c.Bind(updatedNote); err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error bining request body json to storage.Note struct (check logs for more info)")
	}
	updatedNote.LastEdit = time.Now()
	err = s.storage.SaveNote(updatedNote)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "error saving note (check logs for more info)")
	}
	return c.NoContent(http.StatusOK)

}
