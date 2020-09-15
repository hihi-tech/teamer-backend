package controller

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"teamer/model"
)

func (ct Controller) AuthMiddleware() []echo.MiddlewareFunc {
	return []echo.MiddlewareFunc{
		middleware.JWTWithConfig(middleware.JWTConfig{
			SigningKey: []byte(JWTSecret),
			ContextKey: "jwt",
		}),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				user := c.Get("jwt").(*jwt.Token)
				claims := user.Claims.(jwt.MapClaims)

				if !claims.VerifyAudience("service", true) {
					return echo.NewHTTPError(http.StatusUnauthorized, "must use service token")
				}

				email := claims["sub"].(string)

				var findUser model.User
				if err := ct.db.Preload("Schools").Preload("Tags").First(&findUser, &model.User{Email: email}).Error; err != nil {
					return echo.NewHTTPError(http.StatusUnauthorized, "unknown user")
				}

				c.Set("user", &findUser)
				return next(c)
			}
		},
	}
}
