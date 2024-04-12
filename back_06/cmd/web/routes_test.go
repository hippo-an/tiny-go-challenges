package main

import (
	"github.com/dev-hippo-an/tiny-go-challenges/back_06/internal/config"
	"github.com/gorilla/mux"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	m := route(&app)

	switch m.(type) {
	case *mux.Router:
	default:
		t.Errorf("type is not mux.Router, type is %T", m)
	}
}
