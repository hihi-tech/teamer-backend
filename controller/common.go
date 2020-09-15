package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"teamer/model"
)

const JWTSecret = "ilsFeI4h2WZbHJLUjMNGQSXtAbYkgarf"

var DefaultBadRequestResponse = echo.NewHTTPError(http.StatusBadRequest, "bad request: form fields either mismatch or in invalid format")

func getUserFromContext(c echo.Context) *model.User {
	return c.Get("user").(*model.User)
}