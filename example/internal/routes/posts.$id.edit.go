package routes

func (r *Routes) PostsEdit(c echo.Context) error {
    return views.PostsEdit().Render(c.Request().Context(), c.Response().Writer)
}
