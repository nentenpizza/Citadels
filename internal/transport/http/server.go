package httpapi

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/citadels/internal/service"
)

type App struct {
	UsersService service.Users
}

func (a *App) Run(addr string) error {
	e := newEcho()

	a.setupRoutes(e)

	return e.Start(addr)
}

func (a *App) setupRoutes(e *echo.Echo) {
	authHandler := NewAuthHandler(a.UsersService)

	api := e.Group("/api")
	v1 := api.Group("/v1")

	v1.POST("/register", authHandler.OnRegister)

}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	return e
}
