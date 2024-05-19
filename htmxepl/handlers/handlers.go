package handlers

import (
	"github.com/dev-hippo-an/tiny-go-challenges/htmxepl/views"
	"net/http"
)

func HandleFoo(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, views.Index())
}
