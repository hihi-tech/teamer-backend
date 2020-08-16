package main

import (
	"gopkg.in/go-playground/validator.v9"
	"time"
)

type Validator struct {
	validator *validator.Validate
}

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type Config struct {
	Server struct {
		Address string `yaml:"address"`
		Hostname string `yaml:"hostname"`
		CORS    struct {
			Enabled      bool     `yaml:"enabled"`
			AllowOrigins []string `yaml:"allowOrigins"`
		} `yaml:"cors"`
	} `yaml:"server"`
	Database struct {
		DSN string `json:"dsn"`
	} `yaml:"database"`
	Email struct {
		SMTP struct {
			Hostname string `yaml:"hostname"`
			Port int `yaml:"port"`
			Username string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"smtp"`
	} `yaml:"email"`
}

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

