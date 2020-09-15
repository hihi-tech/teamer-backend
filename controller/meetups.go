package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

type Meetup struct {
	gorm.Model

	Name string `json:"name" `
	Description string `gorm:"type:text" json:"description"`
	Location string `json:"location"`

	Start *time.Time `json:"start"`
	End *time.Time `json:"end"`
	Members []User `json:"members"`
	Tags []Tag `json:"tags"`
}

type CreateMeetupRequestForm struct {
	Name string `json:"name" validate:"max=64,required"`
	Description string `gorm:"type:text" json:"description" validate:"max=2048"`
	Location string `json:"location"`

	Start *time.Time `json:"start"`
	End *time.Time `json:"end"`
	Members []uint `json:"members"`
}

func meetupCreateMeetup(c echo.Context) error {
	var form CreateMeetupRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	toSave := Meetup{
		Name: form.Name,
		Description: form.Description,
		Location: form.Location,
		Start: form.Start,
		End: form.End,
	}

	var users []*User
	for _, member := range form.Members {
		var foundUser User
		if err := DB.First(&foundUser, member).Error; err != nil {
			spew.Dump(err)
			return echo.NewHTTPError(http.StatusBadRequest, "cannot found member with id " + strconv.Itoa(int(member)))
		}
		users = append(users, &foundUser)
	}

	if err := DB.Create(&toSave).Error; err != nil {
		LogDb.Println("create meetup: failed to create db record: " + spew.Sdump(err))
		//return echo.NewHTTPError(http.StatusInternalServerError, "failed to create db record")
		return fmt.Errorf("failed to create db record")
	}

	return c.JSON(http.StatusOK, toSave)
}
