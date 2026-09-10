package main

import (
	"fmt"
	"net/http"

	"carbonite/client/internal/handler"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) routes() http.Handler {
	e := echo.New()
	e.HTTPErrorHandler = customHTTPErrorHandler

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(app.store))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:_csrf",
		CookieSecure:   true,
		CookieHTTPOnly: true,
	}))

	// Navigation
	homeHandler := handler.HomeHandler{
		Store:       app.store,
		SessionName: app.sessionName,
	}
	e.GET("/", homeHandler.HandleHomepage)
	e.File("/favicon.ico", "images/favicon.ico")

	// Welcome
	welcomeHandler := handler.WelcomeHandler{}
	gw := e.Group("/welcome")
	gw.Use(app.Auth)
	gw.GET("", welcomeHandler.HandleWelcomepage)

	// Login
	loginHandler := handler.LoginHandler{
		Store:       app.store,
		SessionName: app.sessionName,
		Requester:   app.requester,
		Cacher:      app.cacher,
	}
	e.GET("/login", loginHandler.HandleGetLogin)
	e.POST("/login", loginHandler.HandlePostLogin)
	e.GET("/logout", loginHandler.HandleGetLogout)

	return e.Server.Handler
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	errorPage := fmt.Sprintf("%d.html", code)
	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}
}
