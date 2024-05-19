package handlers

import (
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

type HTTPHandler func(w http.ResponseWriter, r *http.Request) error

// Make is decorate / adapter of HTTPHandler
func Make(h HTTPHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("HTTP Handler error", "err", err, "path", r.URL.Path)
		}
	}
}

func Render(w http.ResponseWriter, r *http.Request, c templ.Component) error {
	return c.Render(r.Context(), w)
}
