package web

import "github.com/quinn/g/example/internal/routes"

func NewServer() *echo.Echo {
	e := echo.New()

	r := &routes.Routes{}

	/* insert new routes here */

	return e
}
