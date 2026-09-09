package handler

import "github.com/labstack/echo/v4"

type (
	ClientHandler struct{}
)

func (h *ClientHandler) HandleGetClients(c echo.Context) error {
	return nil
}

func (h *ClientHandler) HandleGetClient(c echo.Context) error {
	return nil
}
