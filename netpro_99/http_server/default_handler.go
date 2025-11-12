package httpserver

import (
	"fmt"
	"io"
	"net/http"
)

func DefaultHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			r.Close()
		}(r.Body)

		switch r.Method {
		case http.MethodGet:
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Hello! This is a GET Request")
		case http.MethodPost:
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, "Hello! This is a POST Request")
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}
