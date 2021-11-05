package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nentenpizza/citadels/internal/repository"
)

type Server struct {
	Repos *repository.Repositories
}

func (s Server) Start(addr string) error {
	e := echo.New()

	s.setupRoutes(e)

	return e.Start(addr)
}

func (s Server) setupRoutes(e *echo.Echo) {

}

func newEcho() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORS())
	return e
}
