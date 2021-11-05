package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nentenpizza/citadels/internal/service"
)

func (h Handler) OnRegister(c echo.Context) error {
	var form service.UserRegisterForm
	if err := c.Bind(&form); err != nil {
		return err
	}
	if err := c.Validate(&form); err != nil {
		return err
	}
	err := h.Users.Register(form)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, echo.Map{"ok": true})
}
