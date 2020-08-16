package main

import "github.com/labstack/echo"

func getUserFromContext(c echo.Context) *User {
	return c.Get("user").(*User)
}