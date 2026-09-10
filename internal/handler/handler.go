package handler

import (
	"fmt"
	"net/http"

	"carbonite/client/internal/data"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

type (
	// Storer interface
	Storer interface {
		Get(*http.Request, string) (*sessions.Session, error)
		New(*http.Request, string) (*sessions.Session, error)
		Save(*http.Request, http.ResponseWriter, *sessions.Session) error
	}

	// Requester interface
	Requester interface {
		PostToken(string, string) (*data.TokenPayload, error)
		PostRefreshToken(string) (*data.TokenPayload, error)
	}
)

// Render replaces Echo's echo.Context.Render() with templ's templ.Component.Render().
func Render(c echo.Context, statusCode int, t templ.Component) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	// if c.Request().Header.Get("Hx-Request") == "true" {
	if err := t.Render(c.Request().Context(), c.Response().Writer); err != nil {
		return c.String(500, "Failed to render webpage.")
	}
	return nil
}

func return500(c echo.Context, e error) error {
	code := http.StatusInternalServerError
	errorPage := fmt.Sprintf("%d.html", code)
	c.Logger().Error(e)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
	return nil
}
