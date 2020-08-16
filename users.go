package main

import (
	"database/sql"
	"github.com/labstack/echo"
	"net/http"
)

// User describes a user object
type User struct {
	DBSafeModel

	Email    string          `json:"email" gorm:"unique" sql:"index"`
	Password string          `json:"-"`
	Phone    *sql.NullString `json:"phone"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Birthday  string `json:"birthday"`

	Schools []*School `json:"school" gorm:"many2many:users_schools"`
	Tags    []*Tag    `json:"tags" gorm:"many2many:users_tags"`
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
