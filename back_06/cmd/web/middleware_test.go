package main

import (
	"net/http"
	"testing"
)

func TestLoggingMiddleware(t *testing.T) {
	var h myHandler

	hd := LoggingMiddleware(&h)

	switch hd.(type) {
	case http.Handler:
	default:
		t.Error("type is not http.Handler")
	}
}

func TestNoSurf(t *testing.T) {
	var h myHandler

	hd := CSRFMiddleware(&h)

	switch hd.(type) {
	case http.Handler:
	default:
		t.Error("type is not http.Handler")
	}
}

func TestSessionLoad(t *testing.T) {
	var h myHandler

	hd := CSRFMiddleware(&h)

	switch hd.(type) {
	case http.Handler:
	default:
		t.Error("type is not http.Handler")
	}
}
