package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r *Routes) RegisterItemsRoutes(g *echo.Group) {
	g.GET("/", r.getItemsList)
}

func (r *Routes) getItemsList(c echo.Context) error {
	r.Log.Info("Get items list")
	return c.String(http.StatusOK, "List of items")
}
