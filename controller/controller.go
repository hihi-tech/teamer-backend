package controllers

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	"log"
)

type Controller struct {
	LogDb *log.Logger
	LogService *log.Logger
	LogWeChat *log.Logger

	DB *gorm.DB
	MailDialer *gomail.Dialer
	Conf Config
}

func NewController() *Controller {
	return &Controller{}
}
