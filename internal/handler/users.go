package handler

import (
	"fmt"
	"net/http"

	"carbonite/client/internal/data"

	"github.com/labstack/echo/v4"
)

type (
	UserHandler struct{}
)

func (h *UserHandler) HandleGetUsers(c echo.Context) error {
	_, logged := c.Get("token").(*data.TokenPayload)

	return c.String(http.StatusOK, fmt.Sprintf("logged = %v ", logged))
}

func (h *UserHandler) HandleGetUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func (h *UserHandler) HandlePostUser(c echo.Context) error {
	return nil
}

func (h *UserHandler) HandlePutUser(c echo.Context) error {
	return nil
}

func (h *UserHandler) HandlePatchUser(c echo.Context) error {
	return nil
}

func (h *UserHandler) HandleDeleteUser(c echo.Context) error {
	return nil
}
