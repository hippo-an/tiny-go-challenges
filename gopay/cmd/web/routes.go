package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *application) routes() http.Handler {
	mux := echo.New()
	return mux
}
