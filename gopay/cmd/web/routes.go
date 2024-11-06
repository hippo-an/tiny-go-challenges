package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (a *application) routes() http.Handler {
	mux := echo.New()
	mux.GET("/virtual-terminal", a.VirtualTerminal)
	return mux
}
