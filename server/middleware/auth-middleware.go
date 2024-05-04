package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/services"
)

func AuthMiddleware(logger *log.Logger, cfg *config.Config, userService services.UsersServiceInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authToken := c.Request().Header.Get("Auth-Token")

			if authToken == "" {
				logger.Errorln("auth token is empty")
				return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
			}

			claims := &services.Claims{}
			_, err := jwt.ParseWithClaims(authToken, claims, func(token *jwt.Token) (any, error) {
				return []byte(cfg.SECRET_KEY), nil
			})
			if err != nil {
				logger.Errorln("failed to parse token:", err)
				return c.JSON(http.StatusInternalServerError, errors.InternalServerErr)
			}

			user := userService.GetUserByEmail(claims.Email)
			if user == nil {
				logger.Errorln("user not found")
				return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
			}
			c.Set("user", user)

			return next(c)
		}
	}
}
