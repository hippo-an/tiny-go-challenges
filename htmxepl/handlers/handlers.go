package handlers

import (
	"github.com/dev-hippo-an/tiny-go-challenges/htmxepl/views/home"
	"net/http"
	"time"
)

func HandleHome(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, home.Home("Homepage"))
}

func HandleGetTestData(w http.ResponseWriter, r *http.Request) error {
	time.Sleep(3 * time.Second)
	w.Write([]byte("wonderful"))
	return nil
}
