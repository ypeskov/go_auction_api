package middleware

import (
	goerrors "errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"ypeskov/go_hillel_9/internal/config"
	"ypeskov/go_hillel_9/internal/errors"
	"ypeskov/go_hillel_9/internal/log"
	"ypeskov/go_hillel_9/services"
)

func AuthMiddleware(logger *log.Logger, cfg *config.Config, userService services.UsersServiceInterface) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authTokenHeader := c.Request().Header.Get("Auth-Token")

			if authTokenHeader == "" {
				logger.Errorln("auth token is empty")
				return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
			}

			tokenHeaderParts := strings.Split(authTokenHeader, "Bearer ")
			if len(tokenHeaderParts) != 2 {
				logger.Errorln("invalid token format. Expected [Bearer <token>]")
				return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
			}
			token := tokenHeaderParts[1]

			claims := &services.Claims{}
			_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
				return []byte(cfg.SecretKey), nil
			})
			if err != nil {
				if goerrors.Is(err, jwt.ErrTokenExpired) {
					logger.Errorln("token expired")
					return c.JSON(http.StatusUnauthorized, errors.TokenExpiredErr)
				}
				logger.Errorln("failed to parse token", err)
				return c.JSON(http.StatusUnauthorized, errors.UnauthorizedErr)
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
