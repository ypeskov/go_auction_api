package routes

import (
	goerrors "errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/repository/models"
)

func (r *Routes) RegisterItemsRoutes(g *echo.Group) {
	g.GET("/", r.getItemsList)
	g.POST("/", r.createItem)
	g.GET("/:id", r.getItem)
	g.PUT("/:id", r.updateItem)
	g.DELETE("/:id", r.deleteItem)
}

// getItemsList retrieves a list of items.
// It returns a JSON array of item details.
// @summary Get Items List
// @tags Items
// @description Retrieves a list of all available items.
// @accept json
// @produce json
// @success 200 {array} models.Item "List of items"
// @failure 500 {object} errors.Error "Internal server error"
// @router /items/ [get]
func (r *Routes) getItemsList(c echo.Context) error {
	r.Log.Infof("Getting items list ...")

	items, err := r.itemsService.GetItemsList()
	if err != nil {
		r.Log.Error("failed to get items from db", err)
		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to get items from db"))
	}

	return c.JSON(http.StatusOK, &items)
}

// createItem creates a new item based on the provided details.
// It returns the created item details or an error if the creation fails.
// @summary Create Item
// @tags Items
// @description Creates a new item based on the provided request body.
// @accept json
// @produce json
// @param item body models.Item true "Item details"
// @success 201 {object} models.Item "Item created successfully"
// @failure 400 {object} errors.Error "Bad Request: Failed to parse request body or validation failed"
// @router /items/ [post]
func (r *Routes) createItem(c echo.Context) error {
	r.Log.Infof("Creating item ...")

	req := new(models.Item)

	err := c.Bind(req)
	if err != nil {
		r.Log.Error("failed to parse request body", err)
		return c.JSON(http.StatusBadRequest,
			errors.NewError("INCORRECT_REQUEST_BODY", "Failed to parse request body"))
	}

	err = req.Validate()
	if err != nil {
		r.Log.Error("validation failed: ", err)
		return c.JSON(http.StatusBadRequest,
			errors.NewError(errors.ValidationFailedErr.Code, err.Error()))
	}

	item, err := r.itemsService.CreateItem(req)
	if err != nil {
		r.Log.Error("failed to create item", err)
		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to create item"))

	}

	r.Log.Infof("Inserted ID: %d", item.Id)

	return c.JSON(http.StatusCreated, &item)
}

// getItem retrieves an item by its ID.
// It returns a JSON object of the item details if found.
// @summary Get Item
// @tags Items
// @description Retrieves the details of an item by its ID.
// @accept json
// @produce json
// @param id path int true "ID of the item to retrieve"
// @success 200 {object} models.Item "Item retrieved successfully"
// @failure 404 {object} errors.Error "Item not found"
// @router /items/{id} [get]
func (r *Routes) getItem(c echo.Context) error {
	r.Log.Infof("Get item with id: %s", c.Param("id"))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.Log.Error("failed to convert id to int!!!", err)
		return c.JSON(http.StatusBadRequest, errors.NewError("INVALID_ID", "Invalid ID"))
	}

	item, err := r.itemsService.GetItemById(id)
	if err != nil {
		r.Log.Error("failed to get item by id", err)
		return c.JSON(http.StatusNotFound, errors.NewError("ITEM_NOT_FOUND", "Item not found"))
	}

	return c.JSON(http.StatusOK, &item)
}

// updateItem updates an item by its ID based on the provided details.
// It returns the updated item details or an error if the update fails.
// @summary Update Item
// @tags Items
// @description Updates the details of an item based on the provided request body and item ID.
// @accept json
// @produce json
// @param id path int true "ID of the item to update"
// @param item body models.Item true "Updated item details"
// @success 200 {object} models.Item "Item updated successfully"
// @failure 400 {object} errors.Error "Bad Request: Failed to parse request body or validation failed"
// @router /items/{id} [put]
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

	item, err := r.itemsService.UpdateItem(id, req)
	if err != nil {
		r.Log.Error("failed to update item", err)
		return c.JSON(http.StatusInternalServerError, errors.NewError("INTERNAL_SERVER_ERROR", "Failed to update item"))
	}

	return c.JSON(http.StatusOK, &item)
}

// deleteItem deletes an item by its ID.
// It returns no content if the deletion is successful, or an error otherwise.
// @summary Delete Item
// @tags Items
// @description Deletes an item by its ID.
// @accept json
// @produce json
// @param id path int true "ID of the item to delete"
// @success 204 "Item deleted successfully"
// @failure 400 {object} errors.Error "Bad Request: Invalid ID"
// @router /items/{id} [delete]
func (r *Routes) deleteItem(c echo.Context) error {
	r.Log.Infof("Delete item with id: %s", c.Param("id"))

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		r.Log.Error("failed to convert id to int!!!", err)
		return c.JSON(http.StatusBadRequest, errors.NewError("INVALID_ID", "Invalid ID"))
	}

	err = r.itemsService.DeleteItem(id)
	if err != nil {
		if goerrors.Is(err, errors.NotFoundErr) {
			r.Log.Error("item not found", err)
			return c.JSON(http.StatusNotFound, errors.NewError("ITEM_NOT_FOUND", "Item not found"))
		}
		r.Log.Errorln("failed to delete item", err)
		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to delete item"))
	}

	return c.NoContent(http.StatusNoContent)
}
