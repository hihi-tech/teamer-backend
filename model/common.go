package model

import "time"

type DBFullModel struct {
	ID        uint `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type DBSafeModel struct {
	ID        uint `json:"-" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

type DBSafeModelWithTime struct {
	ID        uint `json:"-" gorm:"primary_key"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}
