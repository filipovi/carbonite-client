package handler

import (
	"net/http"

	"carbonite/client/internal/data"
	"carbonite/client/internal/validator"
	"carbonite/client/templates/pages"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	LoginHandler struct {
		Store       Storer
		SessionName string
		Requester   Requester
		Cacher      data.Cacher
	}
)

func (h *LoginHandler) HandlePostLogin(c echo.Context) error {
	session, _ := h.Store.Get(c.Request(), h.SessionName)
	if _, ok := session.Values["user_id"].(string); ok {
		c.Redirect(http.StatusSeeOther, "/")
	}

	b := new(data.LoginBody)
	if err := c.Bind(b); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	v := validator.New()
	if data.ValidateLoginBody(v, b); !v.Valid() {
		err := v.Errors
		return Render(c, http.StatusOK, pages.Login("Page de connexion", *b, csrfToken, err))
	}

	// Get the token
	token, err := h.Requester.PostToken(b.Email, b.Password)
	if err != nil {
		return return500(c, err)
	}

	if token.AuthenticationToken.Scope != "authentication" {
		err := map[string]string{
			"email":    "mauvais identifiant",
			"password": "mauvais identifiant",
		}

		return Render(c, http.StatusOK, pages.Login("Page de connexion", *b, csrfToken, err))
	}

	key, err := token.Save(h.Cacher)
	if err != nil {
		return return500(c, err)
	}

	session.Values["user_id"] = key
	if err = session.Save(c.Request(), c.Response()); err != nil {
		return return500(c, err)
	}
	c.Redirect(http.StatusSeeOther, "/welcome")
	return nil
}

func (h *LoginHandler) HandleGetLogout(c echo.Context) error {
	var key string
	session, _ := h.Store.Get(c.Request(), h.SessionName)
	key, ok := session.Values["user_id"].(string)
	if !ok {
		c.Redirect(http.StatusSeeOther, "/")
		return nil
	}
	session.Values["user_id"] = nil
	session.Options.MaxAge = -1

	err := data.RemoveToken(h.Cacher, key)
	if err != nil {
		return return500(c, err)
	}

	if err = session.Save(c.Request(), c.Response()); err != nil {
		return return500(c, err)
	}

	c.Redirect(http.StatusSeeOther, "/")
	return nil
}

func (h *LoginHandler) HandleGetLogin(c echo.Context) error {
	sess, _ := h.Store.Get(c.Request(), h.SessionName)
	if _, ok := sess.Values["user_id"].(string); ok {
		c.Redirect(http.StatusSeeOther, "/")
	}
	csrfToken := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	return Render(c, http.StatusOK, pages.Login("Page de connexion", data.LoginBody{}, csrfToken, nil))
}
