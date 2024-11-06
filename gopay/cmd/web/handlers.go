package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (a *application) VirtualTerminal(c echo.Context) error {
	if err := a.renderTemplate(c, "terminal", nil); err != nil {
		a.errorLog.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return nil
}
