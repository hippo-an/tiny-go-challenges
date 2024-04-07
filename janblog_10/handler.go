package main

import (
	"bytes"
	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Post struct {
	Title       string    `toml:"title"`
	Slug        string    `toml:"slug"`
	Description string    `toml:"description"`
	Date        time.Time `toml:"date"`
	Content     template.HTML
	Author      Author `toml:"author"`
}

type Author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

func PostHandler(sl SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post Post
		post.Slug = r.PathValue("slug")
		postMarkdown, err := sl.Read(post.Slug)
		if err != nil {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}

		rest, err := frontmatter.Parse(strings.NewReader(postMarkdown), &post)
		if err != nil {
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
			return
		}

		mdRenderer := goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("dracula")),
			),
		)
		var buf bytes.Buffer

		err = mdRenderer.Convert(rest, &buf)

		if err != nil {
			http.Error(w, "Error converting markdown", http.StatusInternalServerError)
			return
		}

		tpl, err := template.ParseFiles("post.gohtml")
		if err != nil {
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}
		post.Content = template.HTML(buf.String())
		err = tpl.Execute(w, post)
	}
}
