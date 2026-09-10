package handler

import (
	"net/http"

	"carbonite/client/templates/pages"

	"github.com/labstack/echo/v4"
)

type (
	HomeHandler struct {
		Store       Storer
		SessionName string
	}
)

func (h *HomeHandler) HandleHomepage(c echo.Context) error {
	session, _ := h.Store.Get(c.Request(), h.SessionName)
	_, logged := session.Values["login_id"].(string)
	if logged {
		c.Redirect(http.StatusSeeOther, "/welcome")
	}
	return Render(c, http.StatusOK, pages.Home("Homepage", logged, "Bienvenue chez Planet55"))
}
