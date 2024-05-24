package routes

import (
	goerrors "errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/repository/models"
	"ypeskov/go_hillel_9/services"
)

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type TokenResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func (r *Routes) RegisterUsersRoutes(g *echo.Group) {
	g.GET("/", r.getUsersList)
	g.POST("/", r.createUser)
	g.POST("/login/", r.LoginUser)
	g.POST("/refresh/", r.getNewAccessToken)
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

	users, err := r.UsersService.GetUsersList()
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

	newUser, err := r.UsersService.CreateUser(req)
	if err != nil {
		r.Log.Error("failed to create user", err)

		return c.JSON(http.StatusInternalServerError,
			errors.NewError("INTERNAL_SERVER_ERROR", "Failed to create user"))
	}
	newUser.PasswordHash = ""

	return c.JSON(http.StatusCreated, newUser)
}

// LoginUser logs in a user based on the provided credentials.
// It returns a JWT token if the login is successful or an error if it fails.
// @summary Login User
// @tags Users
// @description Logs in a user based on the provided credentials.
// @accept json
// @produce json
// @param user body Credentials true "User credentials"
// @success 200 {object} string "JWT"
// @failure 400 {object} errors.Error "Bad Request"
// @failure 401 {object} errors.Error "Unauthorized"
// @router /users/login/ [post]
func (r *Routes) LoginUser(c echo.Context) error {
	r.Log.Infof("Logging in user ...")

	var creds Credentials
	err := c.Bind(&creds)
	if err != nil {
		r.Log.Errorln("failed to parse request body", err)

		return c.JSON(http.StatusBadRequest, errors.BadRequestErr)
	}

	minutes := time.Duration(r.cfg.AccessTokenLifetimeMinutes)
	token, err := r.UsersService.GetJWT(creds.Email, creds.Password, minutes, false)
	if err != nil {
		r.Log.Errorln("failed to get JWT", err)

		return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
	}

	refreshToken, err := r.UsersService.GetRefreshToken(creds.Email, creds.Password, false)
	if err != nil {
		r.Log.Errorln("failed to get refresh token", err)

		return c.JSON(http.StatusInternalServerError, errors.InternalServerErr)
	}

	return c.JSON(http.StatusOK, &TokenResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

func (r *Routes) getNewAccessToken(c echo.Context) error {
	r.Log.Infof("Refreshing access token ...")
	refreshToken := c.Request().Header.Get("Refresh-Token")

	// Check if the refresh token is not expired
	claims := &services.Claims{}
	_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (any, error) {
		return []byte(r.cfg.SecretKey), nil
	})
	if err != nil {
		if goerrors.Is(err, jwt.ErrTokenExpired) {
			r.Log.Errorln("token expired")

			return c.JSON(http.StatusUnauthorized, errors.TokenExpiredErr)
		}
		r.Log.Errorln("failed to parse token", err)

		return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
	}

	// then check if refresh token is valid and not used
	user, err := r.UsersService.GetUserByRefreshToken(refreshToken)
	if err != nil {
		r.Log.Errorln("failed to get user by refresh accessToken", err)

		return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
	}

	// check if user in DB is the same as in token
	if user.Email != claims.Email && user.Id != claims.Id {
		r.Log.Errorln("user in token is not the same as in DB")

		return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
	}

	minutes := time.Duration(r.cfg.AccessTokenLifetimeMinutes)
	accessToken, err := r.UsersService.GetJWT(user.Email, "", minutes, true)
	if err != nil {
		r.Log.Errorln("failed to get JWT", err)

		return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
	}

	newRefreshToken, err := r.UsersService.GetRefreshToken(user.Email, "", true)
	if err != nil {
		r.Log.Errorln("failed to get refresh token", err)

		return c.JSON(http.StatusInternalServerError, errors.InternalServerErr)
	}

	return c.JSON(http.StatusOK, &TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	})
}
