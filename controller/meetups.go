package controller

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"teamer/model"
	"time"
)

type CreateMeetupRequestForm struct {
	Name string `json:"name" validate:"max=64,required"`
	Description string `gorm:"type:text" json:"description" validate:"max=2048"`
	Location string `json:"location"`

	Start *time.Time `json:"start"`
	End *time.Time `json:"end"`
	Members []uint `json:"members"`
}

// CreateMeetup godoc
// @Summary Create a meetup
// @Description Create a meetup
// @Tags Meetup
// @Accept json
// @Produce json
// @Param body body CreateMeetupRequestForm true "Request body"
// @Success 200 {object} model.Meetup
// @Failure 400 {object} model.HTTPError
// @Failure 404 {object} model.HTTPError
// @Failure 500 {object} model.HTTPError
// @Security JwtAuth
// @Router /meetups [put]
func (ct Controller) CreateMeetup(c echo.Context) error {
	var form CreateMeetupRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	toSave := model.Meetup{
		Name: form.Name,
		Description: form.Description,
		Location: form.Location,
		Start: form.Start,
		End: form.End,
	}

	var users []*model.User
	for _, member := range form.Members {
		var foundUser model.User
		if err := ct.db.First(&foundUser, member).Error; err != nil {
			spew.Dump(err)
			return echo.NewHTTPError(http.StatusBadRequest, "cannot found member with id " + strconv.Itoa(int(member)))
		}
		users = append(users, &foundUser)
	}

	if err := ct.db.Create(&toSave).Error; err != nil {
		ct.logger.Println("create meetup: failed to create db record: " + spew.Sdump(err))
		//return echo.NewHTTPError(http.StatusInternalServerError, "failed to create db record")
		return fmt.Errorf("failed to create db record")
	}

	return c.JSON(http.StatusOK, toSave)
}
