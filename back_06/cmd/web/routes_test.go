package main

import (
	"testing"

	"github.com/gorilla/mux"
	"github.com/hippo-an/tiny-go-challenges/back_06/internal/config"
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
