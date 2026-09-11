package main

import (
	"net/http"

	"carbonite/client/internal/data"

	"github.com/labstack/echo/v4"
)

// Auth middleware
func (app *application) Auth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := app.store.Get(c.Request(), app.sessionName)
		ID, ok := session.Values["user_id"].(string)
		if !ok {
			c.Logger().Warn("ID is not a string")
			c.Redirect(http.StatusSeeOther, "/login")
			return nil
		}

		token, err := data.LoadToken(app.cacher, ID)
		if err != nil {
			c.Logger().Warn("Token not found")
			c.Redirect(http.StatusSeeOther, "/logout") // To remove the cookie
			return nil
		}

		if !token.IsExpired() {
			c.Set("token", token)
			return next(c)
		}

		t, err := app.requester.PostRefreshToken(token.AuthenticationToken.Refresh)
		if err != nil {
			c.Logger().Warn(err)
			c.Redirect(http.StatusSeeOther, "/logout") // To remove the cookie
			return nil
		}

		if err = t.Replace(app.cacher, ID); err != nil {
			c.Logger().Warn(err)
			c.Redirect(http.StatusSeeOther, "/logout") // To remove the cookie
			return nil
		}

		c.Set("token", t)
		return next(c)
	}
}
