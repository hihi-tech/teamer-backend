package controller

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"net/http"
	"teamer/model"
)

type SchoolAddRequestForm struct {
	Name     string `json:"name" validate:"max=256,required"`
	Location string `json:"location" validate:"max=64,required"`
}

type SchoolSearchRequestForm struct {
	Query     string `query:"q" validate:"max=256,required"`
}

func (ct Controller) SearchSchool(c echo.Context) error {
	var form SchoolSearchRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	var schools []model.School
	ct.db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", form.Query)).Limit(10).Find(&schools)
	return c.JSON(http.StatusOK, schools)
}

func (ct Controller) AddSchool(c echo.Context) error {
	var form SchoolAddRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	if err := ct.db.Create(&model.School{Name: form.Name, Location: form.Location}).Error; err != nil {
		ct.logger.Println("add school: failed to create record on db: " + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create record on db")
	}

	return c.NoContent(http.StatusAccepted)
}
