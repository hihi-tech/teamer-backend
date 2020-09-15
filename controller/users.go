package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"teamer/model"
)

// UserGetAll godoc
// @Summary Get all users
// @Description Get all user information to display in discover page. Notice that this API is subject to change in the future due to security issues of exposing all users directly.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {array} model.User
// @Failure 400 {object} model.HTTPError
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Security JwtAuth
// @Router /users/all [get]
func (ct Controller) UserGetAll(c echo.Context) error {
	var users []model.User
	ct.db.Preload("Schools").Preload("Tags").Limit(10).Find(&users)
	return c.JSON(http.StatusOK, users)
}

// UserGetProfile godoc
// @Summary Get current User profile
// @Description Get current user profile
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 400 {object} model.HTTPError
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Security JwtAuth
// @Router /users/profile [get]
func (ct Controller) UserGetProfile(c echo.Context) error {
	user := getUserFromContext(c)
	return c.JSON(http.StatusOK, &user)
}

// UserPatchProfile godoc
// @Summary Patch current User profile
// @Description Patch current user profile
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} model.User
// @Failure 400 {object} model.HTTPError
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Security JwtAuth
// @Router /users/profile [patch]
func (ct Controller) UserPatchProfile(c echo.Context) error {
	user := getUserFromContext(c)
	user.Schools = []*model.School{
		{Name: "杭州市第一中学", Location: "杭州"},
		{Name: "杭州市第二中学", Location: "杭州"},
	}
	user.Tags = []*model.Tag{
		{Name: "前端开发", Description: "就，就前端开发嘛"},
		{Name: "后端开发", Description: "后端开发嘛诶嘿"},
	}
	if err := ct.db.Save(&user).Error; err != nil {
		ct.logger.Println("patch profile: unable to save user record to db: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to save user record")
	}
	return c.JSON(http.StatusAccepted, &user)
}
