package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/citadels/internal/service"
	v1 "github.com/nentenpizza/citadels/internal/transport/http/v1"
)

type Server struct {
	Users service.Users
}

func (s Server) Start(addr string) error {
	e := newEcho()

	s.setupRoutes(e)

	return e.Start(addr)
}

func (s Server) setupRoutes(e *echo.Echo) {
	handlerV1 := v1.NewHandler(s.Users)

	api := e.Group("/api")
	apiV1 := api.Group("/v1")

	apiV1.POST("/register", handlerV1.OnRegister)

}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	return e
}
