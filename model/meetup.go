package model

import (
	"github.com/jinzhu/gorm"
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
