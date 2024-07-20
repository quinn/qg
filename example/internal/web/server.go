package web

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/g/example/internal/routes"
)

func NewServer() *echo.Echo {
	e := echo.New()

	r := &routes.Routes{}

	/* insert new routes here */

	return e
}
