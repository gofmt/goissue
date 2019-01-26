package api

import (
	"github.com/labstack/echo"
)

func InitRouter(g *echo.Group) {
	g.GET("", func(c echo.Context) error {
		return c.Render(200, "home", nil)
	})
}
