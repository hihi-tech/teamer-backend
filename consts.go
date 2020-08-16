package main

import (
	"github.com/labstack/echo"
	"net/http"
)

const JWTSecret = "ilsFeI4h2WZbHJLUjMNGQSXtAbYkgarf"

var (
	DefaultBadRequestResponse = echo.NewHTTPError(http.StatusBadRequest, "bad request: form fields either mismatch or in invalid format")
)

