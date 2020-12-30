package storage

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

type Note struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Contents string    `json:"contents"`
	LastEdit time.Time `json:"last_edit"`
}

func (note *Note) RenderHTML() string {
	parser := parser.NewWithExtensions(parser.CommonExtensions)
	html := markdown.ToHTML([]byte(note.Contents), parser, nil)
	return string(html)
}

func (note *Note) ParseTitle() string {
	title := ""

	parser := parser.NewWithExtensions(parser.CommonExtensions)
	rootNode := parser.Parse([]byte(note.Contents))

	// return contents of first matched heading
	ast.WalkFunc(rootNode, func(node ast.Node, entering bool) ast.WalkStatus {
		if entering {
			switch n := node.(type) {
			case *ast.Heading:
				title = string(n.Children[0].AsLeaf().Literal)
				return ast.Terminate // terminate on first parsed heading
			}
		}
		return ast.GoToNext // continue with next node
	})
	return title
}

func (note *Note) ParseTags() []string {
	// TODO
	return make([]string, 0)
}

func (note *Note) ParseBlobIDs() []string {

	parser := parser.NewWithExtensions(parser.CommonExtensions)
	rootNode := parser.Parse([]byte(note.Contents))

	// candidateURLs will contain all parsed URLs of all Images and Links. Some of them
	// will probably point to non-user-uploaded blobs, so they will have to be filtered.
	candidateURLs := make([]string, 0)
	ast.WalkFunc(rootNode, func(node ast.Node, entering bool) ast.WalkStatus {
		// traverse all Images and Links (both could contain a user-uploaded blob)
		if entering {
			switch n := node.(type) {
			case *ast.Image:
				candidateURLs = append(candidateURLs, string(n.Destination))
			case *ast.Link:
				candidateURLs = append(candidateURLs, string(n.Destination))
			}
		}
		return ast.GoToNext // continue with next node
	})

	blobIDs := make([]string, 0)
	// go through all candidateURLs. if a certain url is a url to a user-uploaded blob,
	// parse ID of the blob.
	for _, url := range candidateURLs {
		trimmedURL := strings.TrimSpace(url)
		if strings.HasPrefix(trimmedURL, "/api/blob/") {
			// the url points to a user-uplaoded blob.
			// parse the blob_id (url is of format: /api/blob/:blob_id/:browser_helper)
			blobID := strings.Split(trimmedURL, "/")[3]
			blobIDs = append(blobIDs, blobID)
		}
	}
	return blobIDs
}

func (note *Note) GenerateLs(storage Storage) error {
	note.ID = "ls"
	note.Title = "ls"
	note.LastEdit = time.Now()

	// generate contents
	lsAsMarkdown := "# All notes\n"
	lsAsMarkdown += "This note is autogenerated, any user changes to it will be ignored.\n\n"
	lsAsMarkdown += "|title|url|ID|size[B]|last edit|\n"
	lsAsMarkdown += "|-----|---|---|------|---------|\n"
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
		noteLine += "**" + note.Title + "** |"
		noteLine += fmt.Sprintf("[/%s](/%s)", noteID, noteID) + "|"
		noteLine += noteID + "|"
		noteLine += fmt.Sprintf("%d", len(note.Contents)) + "|"
		noteLine += note.LastEdit.String() + "|\n"
		lsAsMarkdown += noteLine
	}
	note.Contents = lsAsMarkdown
	return nil
}

func (note *Note) GenerateLsTag(storage Storage) error {
	note.ID = "lstag"
	note.Title = "lstag"
	note.LastEdit = time.Now()

	// generate contents
	lsTagAsMarkdown := "# All note tags\n"
	lsTagAsMarkdown += "This note is autogenerated, any user changes to it will be ignored.\n\n"
	lsTagAsMarkdown += "|tag|note urls|\n"
	lsTagAsMarkdown += "|---|---------|\n"

	tagIndex, err := storage.GetAllNoteTags()
	if err != nil {
		return err
	}

	allNoteTags := make([]string, 0)
	for tag := range tagIndex {
		allNoteTags = append(allNoteTags, tag)
	}

	sort.Strings(allNoteTags)

	for _, tag := range allNoteTags {
		lsTagAsMarkdown += "|`" + tag + "`|"
		for _, noteID := range tagIndex[tag] {
			lsTagAsMarkdown += fmt.Sprintf("[%s](/%s) ", noteID, noteID)
		}
		lsTagAsMarkdown += "|\n"
	}

	note.Contents = lsTagAsMarkdown
	return nil
}
