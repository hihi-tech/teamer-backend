package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo"
	"net/http"
)

type School struct {
	DBFullModel

	Name string `json:"name" gorm:"unique" sql:"index"`
	//AlternateName string `json:"alternateName"`

	// Location is the administration assignment of the current object. Example: `北京` for `北京大学`
	Location string `json:"location" sql:"index"`
	//Description   string `json:"description" gorm:"type:text"`
}

type SchoolAddRequestForm struct {
	Name     string `json:"name" validate:"max=256,required"`
	Location string `json:"location" validate:"max=64,required"`
}

type SchoolSearchRequestForm struct {
	Query     string `query:"q" validate:"max=256,required"`
}

func schoolSearch(c echo.Context) error {
	var form SchoolSearchRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	var schools []School
	DB.Where("name LIKE ?", fmt.Sprintf("%%%s%%", form.Query)).Limit(10).Find(&schools)
	return c.JSON(http.StatusOK, schools)
}

func schoolAdd(c echo.Context) error {
	var form SchoolAddRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	if err := DB.Create(&School{Name: form.Name, Location: form.Location}).Error; err != nil {
		LogDb.Println("add school: failed to create record on db: " + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create record on db")
	}

	return c.NoContent(http.StatusAccepted)
}
