package main

import "github.com/jinzhu/gorm"

type Meetup struct {
	gorm.Model

	Name string `json:"name" validate:"max=64,required"`
	Description string `gorm:"type:text" json:"description" validate:"max=2048"`
	Location string `json:"location"`
	Time *TimeRange `json:"time"`
	Invitees []User `json:"invitees"`
	Tags []Tag `json:"tags"`
}
