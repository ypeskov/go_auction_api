package routes

import (
	"math/rand"
	"net/http"
	"strconv"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/internal/models"

	"github.com/labstack/echo/v4"
)

func (r *Routes) RegisterItemsRoutes(g *echo.Group) {
	g.GET("/", r.getItemsList)
	g.POST("/", r.createItem)
	g.GET("/:id", r.getItem)
	g.PUT("/:id", r.updateItem)
	g.DELETE("/:id", r.deleteItem)
}

func (r *Routes) getItemsList(c echo.Context) error {
	r.Log.Info("Get items list")
	items := []*models.Item{
		{
			ID:           rand.Intn(999_999) + 1,
			Title:        "Item 1",
			InitialPrice: 100,
			Description:  "Description of item 1",
		},
		{
			ID:           rand.Intn(999_999) + 1,
			Title:        "Item 2",
			InitialPrice: 200,
			Description:  "Description of item 2",
		},
	}

	return c.JSON(http.StatusOK, &items)
}

func (r *Routes) createItem(c echo.Context) error {
	r.Log.Infof("Creating item ...")

	req := new(models.Item)

	err := c.Bind(req)
	if err != nil {
		r.Log.Error("failed to parse request body", err)
		return c.JSON(http.StatusBadRequest, errors.NewError("INCORRECT_REQUEST_BODY", "Failed to parse request body"))
	}

	err = req.Validate()
	if err != nil {
		r.Log.Error("validation failed: ", err)
		return c.JSON(http.StatusBadRequest, errors.NewError(errors.ValidationFailedErr.Code, err.Error()))
	}

	req.ID = rand.Intn(999_999) + 1
	r.Log.Infof("Item created: %+v", req)

	return c.JSON(http.StatusCreated, &req)
}

func (r *Routes) getItem(c echo.Context) error {
	r.Log.Infof("Get item with id: %s", c.Param("id"))

	return c.JSON(http.StatusOK, &models.Item{
		ID:           rand.Intn(999_999) + 1,
		Title:        "Item 1",
		InitialPrice: 100,
		Description:  "Description of item 1",
	})
}

func (r *Routes) updateItem(c echo.Context) error {
	r.Log.Infof("Update item with id: %s", c.Param("id"))

	req := new(models.Item)

	err := c.Bind(req)
	if err != nil {
		r.Log.Error("failed to parse request body", err)
		return c.JSON(http.StatusBadRequest, errors.NewError("INCORRECT_REQUEST_BODY", "Failed to parse request body"))
	}

	err = req.Validate()
	if err != nil {
		r.Log.Error("validation failed: ", err)
		return c.JSON(http.StatusBadRequest, errors.NewError(errors.ValidationFailedErr.Code, err.Error()))
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.Log.Error("failed to convert id to int!!!", err)
		return c.JSON(http.StatusBadRequest, errors.NewError("INVALID_ID", "Invalid ID"))
	}
	req.ID = id

	return c.JSON(http.StatusOK, &req)
}

func (r *Routes) deleteItem(c echo.Context) error {
	r.Log.Infof("Delete item with id: %s", c.Param("id"))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.Log.Error("failed to convert id to int!!!", err)
		return c.JSON(http.StatusBadRequest, errors.NewError("INVALID_ID", "Invalid ID"))
	}
	r.Log.Infof("Item with id %d deleted", id)

	return c.NoContent(http.StatusNoContent)
}
