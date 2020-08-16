package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

var MiddlewareRequireAuth = []echo.MiddlewareFunc{
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

			var findUser User
			if err := DB.Preload("Schools").Preload("Tags").First(&findUser, &User{Email: email}).Error; err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "unknown user")
			}

			c.Set("user", &findUser)
			return next(c)
		}
	},
}


