package handler

import (
	"carbonite/admin/templates/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	WelcomeHandler struct{}
)

func (h *WelcomeHandler) HandleWelcomepage(c echo.Context) error {
	return Render(c, http.StatusOK, pages.Welcome("Welcomepage"))
}
