package storage

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Note struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Contents string    `json:"contents"`
	Tags     []string  `json:"tags"`
	LastEdit time.Time `json:"last_edit"`
	IsPublic bool      `json:"is_public"`
}

func (note *Note) ParseTitle() string {
	// this regex extracts the first title (#) or subtitle(##) or any title.
	// TODO(sbrki): use gomarkdown and extract title from AST.
	regex := regexp.MustCompile(`#[^#]+`)
	result := regex.FindString(note.Contents)
	// if no title is found, return empty string
	if len(result) <= 1 {
		return ""
	}
	// if multiple lines are matched, use only first one
	title := strings.Split(result, "\n")[0]
	// remove the preifx # that regex has matched and remove any whitespace
	title = strings.TrimSpace(strings.Split(title, "#")[1])
	return title
}

func (note *Note) GenerateLs(storage Storage) error {
	note.ID = "ls"
	note.Title = "ls"
	note.LastEdit = time.Now()
	note.Tags = make([]string, 1)
	note.Tags[0] = "snote/system"
	note.IsPublic = false
	// generate contents

	lsAsMarkdown := "# All notes\n"
	lsAsMarkdown += "This note is autogenerated, any user changes to it note will be ignored.\n\n"
	lsAsMarkdown += "|url|ID|size[B]|tags|last edit|public|\n"
	lsAsMarkdown += "|---|---|------|----|---------|------|\n"
	allNoteIDs, err := storage.GetAllNoteIDs()
	if err != nil {
		return err
	}

	for _, noteID := range allNoteIDs {
		note, err := storage.LoadNote(noteID)
		if err != nil {
			return err
		}
		noteLine := "|"
		noteLine += fmt.Sprintf("[/%s](/%s)", noteID, noteID) + "|"
		noteLine += noteID + "|"
		noteLine += fmt.Sprintf("%d", len(note.Contents)) + "|"
		noteLine += strings.Join(note.Tags, ",") + "|"
		noteLine += note.LastEdit.String() + "|"
		noteLine += fmt.Sprintf("%t", note.IsPublic) + "|\n"
		lsAsMarkdown += noteLine
	}
	note.Contents = lsAsMarkdown
	return nil
}
