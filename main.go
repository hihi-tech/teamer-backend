package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/gomail.v2"
	"log"
	"net/http"
	"time"
)

var (
	LogDb *log.Logger
	LogService *log.Logger
	LogWeChat *log.Logger

	DB *gorm.DB
	MailDialer *gomail.Dialer
	Conf Config
)

func main() {
	LogDb = NewLogger("db", "Database")
	LogService = NewLogger("service", "Service")
	LogWeChat = NewLogger("wechat", "WeChat")

	LogService.Println("initializing configuration...")
	// load configurations
	err := configor.Load(&Conf, "config.yml")
	if err != nil {
		LogDb.Panic("configuration file error: ", err)
	}

	LogService.Println("initialized configuration as following:")
	spew.Dump(Conf)

	LogService.Println("connecting to database...")
	// initialize database connection
	if DB, err = gorm.Open("mysql", Conf.Database.DSN); err != nil {
		LogDb.Panic("failed to open database", err)
	}

	// initialize database tables
	DB.AutoMigrate(&User{}, &Tag{}, &School{})

	//DB = DB.Debug()

	LogService.Println("database connected & auto migrated...")

	// initialize email
	MailDialer = gomail.NewDialer(
		Conf.Email.SMTP.Hostname,
		Conf.Email.SMTP.Port,
		Conf.Email.SMTP.Username,
		Conf.Email.SMTP.Password,
	)

	LogService.Println("initializing echo router...")
	//// BEGIN - echo Router Declarations ////

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} | ${status} ${method} ${uri} ${latency_human}\n",
	}))
	if Conf.Server.CORS.Enabled {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			//AllowOrigins: Conf.Server.CORS.AllowOrigins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
			MaxAge: int((time.Hour * 24).Seconds()),
		}))
	}

	e.Validator = &Validator{
		validator: validator.New(),
	}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct {
			Message string `json:"message"`
		}{
			Message: "Hi and welcome to Teamer ;)",
		})
	})

	api := e.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authLoginHandler)
			auth.GET("/verify/email/:key", authVerifyEmailHandler)
			auth.POST("/register", authRegisterHandler)
		}
		users := api.Group("/user")
		{
			users.GET("/ok", func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}, MiddlewareRequireAuth...)
			users.GET("/all", userGetAll)
			users.GET("/profile", userGetProfile, MiddlewareRequireAuth...)
			users.PATCH("/profile", userPatchProfile, MiddlewareRequireAuth...)
		}
		schools := api.Group("/school", MiddlewareRequireAuth...)
		{
			schools.GET("/search", schoolSearch)
			schools.PUT("", schoolAdd)
		}
		meetups := api.Group("/meetups", MiddlewareRequireAuth...)
		{
			meetups.PUT("", meetupCreateMeetup)
		}
	}

	LogService.Println("setup completed. starting up server...")

	LogService.Printf("Registered Routers: %v", spew.Sdump(e.Routes()))

	log.Fatalln(e.Start(Conf.Server.Address))
}