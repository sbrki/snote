package storage

import (
	"regexp"
	"strings"
	"time"
)

type Note struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Contents string    `json:"contents"`
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
