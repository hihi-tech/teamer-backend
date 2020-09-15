package main

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"teamer/model"
)

func userGetAll(c echo.Context) error {
	var users []model.User
	DB.Preload("Schools").Preload("Tags").Limit(10).Find(&users)
	return c.JSON(http.StatusOK, users)
}

func userGetProfile(c echo.Context) error {
	user := getUserFromContext(c)
	return c.JSON(http.StatusOK, &user)
}

func userPatchProfile(c echo.Context) error {
	user := getUserFromContext(c)
	user.Schools = []*School{
		{Name: "杭州市第一中学", Location: "杭州"},
		{Name: "杭州市第二中学", Location: "杭州"},
	}
	user.Tags = []*Tag{
		{Name: "前端开发", Description: "就，就前端开发嘛"},
		{Name: "后端开发", Description: "后端开发嘛诶嘿"},
	}
	if err := DB.Save(&user).Error; err != nil {
		LogService.Println("patch profile: unable to save user record to db: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to save user record")
	}
	return c.JSON(http.StatusAccepted, &user)
}