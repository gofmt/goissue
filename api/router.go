package api

import (
	"github.com/labstack/echo"
)

func InitRouter(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.GET("/login", func(c echo.Context) error {
		return c.Render(200, "login", nil)
	})

	e.GET("/register", func(c echo.Context) error {
		return c.Render(200, "register", nil)
	})

	e.GET("/publish", func(c echo.Context) error {
		return c.Render(200, "publish", nil)
	})
}
