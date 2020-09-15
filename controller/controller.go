package controller

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/gomail.v2"
	"io"
	"log"
	"os"
	"path"
	"teamer/model"
)

type Controller struct {
	config model.Config

	db         *gorm.DB
	mailDialer *gomail.Dialer
	logger     *log.Logger
}

const logDirectory = "logs/"

func NewLogger(slug string, name string) *log.Logger {
	filename := path.Join(logDirectory, fmt.Sprintf("%s.log", slug))
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	stream := io.MultiWriter(os.Stdout, logFile)

	return log.New(stream, name, log.Ldate|log.Ltime|log.Lshortfile)
}

func NewController(config model.Config) *Controller {
	Log := NewLogger("service", "Service")

	Log.Println("connecting to database...")
	// initialize database connection
	DB, err := gorm.Open("mysql", config.Database.DSN)
	if err != nil {
		Log.Panic("failed to open database", err)
	}

	DB.AutoMigrate(&model.User{}, &model.Tag{}, &model.School{})

	mailDialer := gomail.NewDialer(
		config.Email.SMTP.Hostname,
		config.Email.SMTP.Port,
		config.Email.SMTP.Username,
		config.Email.SMTP.Password,
	)

	return &Controller{
		config: config,

		db:         DB,
		mailDialer: mailDialer,
		logger:     Log,
	}
}
