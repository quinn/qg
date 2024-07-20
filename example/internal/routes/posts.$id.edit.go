package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/quinn/g/example/views"
)

func (r *Routes) PostsEdit(c echo.Context) error {
	return views.PostsEdit().Render(c.Request().Context(), c.Response().Writer)
}
