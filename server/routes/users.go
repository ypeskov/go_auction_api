package routes

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/repository/models"
)

func (r *Routes) RegisterUsersRoutes(g *echo.Group) {
	g.GET("/", r.getUsersList)
	g.POST("/", r.createUser)
	//g.GET("/:id", r.getItem)
	//g.PUT("/:id", r.updateItem)
	//g.DELETE("/:id", r.deleteItem)
}

// getUsersList retrieves a list of users.
// It returns a JSON array of user details.
// @summary Get Users List
// @tags Users
// @description Retrieves a list of all available users.
// @accept json
// @produce json
// @success 200 {array} models.User "List of users"
// @failure 500 {object} errors.Error "Internal server error"
// @router /users/ [get]
func (r *Routes) getUsersList(c echo.Context) error {
	r.Log.Infof("Getting users list ...")

	users, err := r.UserRepo.GetUsersList()
	if err != nil {
		r.Log.Error("failed to get users from db", err)
		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to get users from db"))
	}

	return c.JSON(http.StatusOK, &users)
}

// createUser creates a new user based on the provided details.
// It returns the created user details or an error if the creation fails.
// @summary Create User
// @tags Users
// @description Creates a new user based on the provided request body.
// @accept json
// @produce json
// @param user body models.User true "User details"
// @success 201 {object} models.User "User created successfully"
// @failure 400 {object} errors.Error "Bad Request: Failed to parse request body or validation failed"
// @router /users/ [post]
func (r *Routes) createUser(c echo.Context) error {
	r.Log.Infof("Creating user ...")

	req := new(models.User)

	err := c.Bind(req)
	if err != nil {
		r.Log.Error("failed to parse request body", err)
		return c.JSON(http.StatusBadRequest,
			errors.NewError("BAD_REQUEST", "Failed to parse request body"))
	}

	err = req.Validate()
	if err != nil {
		r.Log.Errorln("failed to validate request body", err)
		return c.JSON(http.StatusBadRequest,
			errors.NewError("VALIDATION_FAILED", err.Error()))
	}

	newUser, err := r.UserRepo.CreateUser(req)
	if err != nil {
		r.Log.Error("failed to create user", err)
		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to create user"))
	}

	return c.JSON(http.StatusCreated, newUser)
}
