package handler

import (
	"net/http"

	"carbonite/client/templates/pages"

	"github.com/labstack/echo/v4"
)

type (
	WelcomeHandler struct{}
)

func (h *WelcomeHandler) HandleWelcomepage(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Welcome("Welcomepage"))
}
