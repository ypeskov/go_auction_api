package routes

import (
	goerrors "errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/services"
)

func (r *Routes) RegisterItemsRoutes(g *echo.Group) {
	g.GET("/", r.getItemsList)
	g.GET("/all", r.getAllItems)
	g.POST("/comments", r.createItemComment)
	g.POST("/", r.createItem)
	g.GET("/:id", r.getItem)
	g.PUT("/:id", r.updateItem)
	g.DELETE("/:id", r.deleteItem)
	g.POST("/attach", r.attachFile)
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
	user := c.Get("user").(*models.User)
	if user == nil {
		r.Log.Error("failed to get user from context")

		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to get user from context"))
	}
	items, err := r.ItemsService.GetItemsList(user.Id)
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

	user := c.Get("user").(*models.User)
	req.UserId = user.Id
	item, err := r.ItemsService.CreateItem(req, user)
	if err != nil {
		r.Log.Errorln("failed to create item", err)
		if goerrors.Is(err, services.IncorrectUserRoleErr) {

			return c.JSON(http.StatusForbidden,
				errors.NewError("INCORRECT_USER_ROLE", "User is not a seller"))
		}

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
		r.Log.Errorln("failed to convert id to int!!!", err)

		return c.JSON(http.StatusBadRequest, errors.NewError("INVALID_ID", "Invalid ID"))
	}

	user := c.Get("user").(*models.User)
	if user == nil {
		r.Log.Error("failed to get user from context")

		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to get user from context"))
	}

	item, err := r.ItemsService.GetItemById(id, user.Id)
	if err != nil {
		r.Log.Errorln("failed to get item by id", err)

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

	user := c.Get("user").(*models.User)
	if user == nil {
		r.Log.Error("failed to get user from context")

		return c.JSON(http.StatusInternalServerError, errors.NewError("INTERNAL_SERVER_ERROR",
			"Failed to get user from context"))
	}

	item, err := r.ItemsService.UpdateItem(id, req, user.Id)
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
		r.Log.Errorln("failed to convert id to int", err)

		return c.JSON(http.StatusBadRequest, errors.NewError("INVALID_ID", "Invalid ID"))
	}

	user := c.Get("user").(*models.User)
	if user == nil {
		r.Log.Error("failed to get user from context")

		return c.JSON(http.StatusInternalServerError, errors.NewError("INTERNAL_SERVER_ERROR",
			"Failed to get user from context"))
	}
	err = r.ItemsService.DeleteItem(id, user.Id)
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

func (r *Routes) getAllItems(c echo.Context) error {
	r.Log.Infof("Getting all items ...")
	items, err := r.ItemsService.GetAllItems()
	if err != nil {
		r.Log.Error("failed to get items from db", err)

		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to get items from db"))
	}

	return c.JSON(http.StatusOK, &items)
}

func (r *Routes) createItemComment(c echo.Context) error {
	r.Log.Infof("Creating item comment ...")

	itemComment := new(models.ItemComment)

	err := c.Bind(itemComment)
	if err != nil {
		r.Log.Error("failed to parse request body", err)

		return c.JSON(http.StatusBadRequest,
			errors.NewError("INCORRECT_REQUEST_BODY", "Failed to parse request body"))
	}

	err = itemComment.Validate()
	if err != nil {
		r.Log.Error("validation failed: ", err)

		return c.JSON(http.StatusBadRequest,
			errors.NewError(errors.ValidationFailedErr.Code, err.Error()))
	}

	user := c.Get("user").(*models.User)
	itemComment.UserId = user.Id
	comment, err := r.ItemsService.CreateItemComment(itemComment)
	if err != nil {
		r.Log.Errorln("failed to create item comment", err)

		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to create item comment"))
	}

	r.Log.Infof("Inserted ID: %d", comment.Id)

	return c.JSON(http.StatusCreated, &comment)
}

func (r *Routes) attachFile(c echo.Context) error {
	r.Log.Infof("Attaching file ...")

	itemId := c.FormValue("itemId")
	if itemId == "" {
		r.Log.Error("itemId is empty")

		return c.JSON(http.StatusBadRequest,
			errors.NewError("INCORRECT_REQUEST_BODY", "Item ID is empty"))
	}
	id, err := strconv.Atoi(itemId)
	if err != nil {
		r.Log.Error("failed to convert id to int!!!", err)

		return c.JSON(http.StatusBadRequest,
			errors.NewError("INVALID_ID", "Invalid ID"))
	}

	file, err := c.FormFile("file")
	if err != nil {
		r.Log.Error("failed to get file from form", err)

		return c.JSON(http.StatusBadRequest,
			errors.NewError("INCORRECT_REQUEST_BODY", "Failed to get file from form"))
	}
	fileName, err := r.ItemsService.AttachFileToItem(id, file)
	if err != nil {
		r.Log.Error("failed to attach file to item", err)

		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to attach file to item"))
	}

	r.Log.Infof("File [%s] attached to item with id: [%d]", *fileName, id)

	return c.JSON(http.StatusOK, "File attached successfully")
}
