package server

import (
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/patrickmn/go-cache"
	"github.com/sbrki/snote/internal/storage"
	"github.com/sbrki/snote/internal/util"
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

	// use the default echo json logger
	s.echo.Use(middleware.Logger())
	s.echo.Logger.SetLevel(log.INFO)

	s.echo.Static("/static", "web/static")
	s.echo.File("/favicon.ico", "web/static/favicon.ico")

	s.setupRoutes()

	return s
}

func (s *Server) setupRoutes() {
	// setup non-api (HTML) handlers
	s.echo.GET("/", s.htmlIndexHandler)
	s.echo.GET("/:note_id", s.htmlNoteHandler)
	s.echo.GET("/:note_id/edit", s.htmlNoteEditHandler)
	// setup api handlers
	// note endpoints
	s.echo.GET("/api/note/:note_id", s.noteGetHandler)
	s.echo.PUT("/api/note/:note_id", s.notePutHandler)
	s.echo.DELETE("/api/note/:note_id", s.noteDeleteHandler)
	s.echo.POST("/api/note", s.noteCollectionPostHandler)
	// blob endpoints
	s.echo.POST("/api/blob", s.blobCollectionPostHandler)
	s.echo.GET("/api/blob/:blob_id/:browser_filename", s.blobGetHandler)

}

// parses all user-uploaded blobs from all notes and deletes blobs
// from the storage if they are not referenced in any note.
func (s *Server) deleteUnusedBlobs() {
	usedBlobIDs := make([]string, 0)
	allNoteIDs, err := s.storage.GetAllNoteIDs()
	if err != nil {
		s.echo.Logger.Error(err)
		return
	}

	// collect all user-uploaded blob IDs from all notes in storage
	for _, noteID := range allNoteIDs {
		note, err := s.storage.LoadNote(noteID)
		if err != nil {
			s.echo.Logger.Error(err)
			return
		}
		noteBlobIDs := note.ParseBlobIDs()
		usedBlobIDs = append(usedBlobIDs, noteBlobIDs...)
	}

	// if there are no user-uploaded blobs in any of the notes, return.
	// also acts as a sanity check, if for some reason note loading/parsing
	// failed (although all fail cases seem to be covered above), prevent
	// deleting all of the blobs currently stored in storage.
	if len(usedBlobIDs) == 0 {
		return
	}

	// collect all blob IDs across all blobs currently stored in storage
	allStorageBlobIDs, err := s.storage.GetAllBlobIDs()
	if err != nil {
		s.echo.Logger.Error(err)
		return
	}

	// delete blobs that are not used, ie. that are present in allStorageBlobIDs
	// but not present in usedBlobIDs.
	for _, storageBlobID := range allStorageBlobIDs {
		used := util.SliceContainsString(usedBlobIDs, storageBlobID)
		if !used {
			err = s.storage.DeleteBlob(storageBlobID)
			if err != nil {
				s.echo.Logger.Error(err)
				return
			}
			s.echo.Logger.Info("deleted blob:" + storageBlobID)
		}
	}

}

// ment to be run as a separate goroutine and do housekeeping tasks.
func (s *Server) backgroundJobs() {
	for {
		time.Sleep(1 * time.Hour)
		s.deleteUnusedBlobs()
	}
}

func (s *Server) Run() {
	// start background jobs
	go s.backgroundJobs()
	s.echo.Logger.Fatal(s.echo.Start(":8081"))
}
