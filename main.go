package main

import (
	"github.com/jinzhu/configor"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"teamer/controller"
	_ "teamer/docs"
	"teamer/model"
	"time"
)

var Conf model.Config

// @title Teamer API
// @version 0.0.1-alpha.1
// @description This is the Teamer API Documentation. You can found contact information regards to the developer of this API and its corresponding documentation below. Notice that this API Documentation is being generated from the actual backend code by using [Swag](https://github.com/swaggo/swag) and its conventional comment annotation on the service implementation code. Due to such reason, there will be a small chance where inconsistencies exist in-between the API Documentation and the actual behavior of the code. The backend development team will strive to keep the API Documentation updated and accurate as possible. This notice just acts as a reminder ;)

// @contact.name Galvin Gao
// @contact.email me@galvingao.com

// @host teamer.localhost
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization
func main() {
	// load configurations
	err := configor.Load(&Conf, "config.yml")
	if err != nil {
		log.Panicf("configuration file error: %v", err)
	}

	ct := controller.NewController(Conf)

	MiddlewareRequireAuth := ct.AuthMiddleware()

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
			auth.POST("/login", ct.Login)
			auth.GET("/verify/email/:key", ct.VerifyEmail)
			auth.POST("/register", ct.Register)
		}
		users := api.Group("/users")
		{
			users.GET("/ok", func(c echo.Context) error {
				return c.NoContent(http.StatusOK)
			}, MiddlewareRequireAuth...)
			users.GET("/all", ct.UserGetAll)
			users.GET("/profile", ct.UserGetProfile, MiddlewareRequireAuth...)
			users.PATCH("/profile", ct.UserPatchProfile, MiddlewareRequireAuth...)
		}
		schools := api.Group("/schools", MiddlewareRequireAuth...)
		{
			schools.GET("/search", ct.SearchSchool)
			schools.PUT("", ct.AddSchool)
		}
		meetups := api.Group("/meetups", MiddlewareRequireAuth...)
		{
			meetups.PUT("", ct.CreateMeetup)
		}
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	log.Println("setup completed. starting up server...")

	log.Fatalln(e.Start(Conf.Server.Address))
}