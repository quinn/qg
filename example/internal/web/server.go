package web

import (
	"github.com/labstack/echo"
	"github.com/quinn/g/example/internal/routes"
)

func NewServer() *echo.Echo {
	e := echo.New()

	r := &routes.Routes{}

	e.GET("/posts/:id/edit", r.PostsEdit)
	/* insert new routes here */

	return e
}
