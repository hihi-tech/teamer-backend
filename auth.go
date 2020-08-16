package main

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type UserLoginRequestForm struct {
	Email string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
	Details *User `json:"details"`
}

type UserRegisterRequestForm struct {
	Email string `json:"email" validate:"email,max=64,required"`
	Password string `json:"password" validate:"min=8,max=64,required"`
	Phone    *string `json:"phone" validate:"max=18,required"`

	FirstName string `json:"firstName" validate:"max=64,required"`
	LastName  string `json:"lastName" validate:"max=64,required"`
	Birthday  string `json:"birthday" validate:"max=24,required"`
}

func createJwt(user *User) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["aud"] = "service"

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func authLoginHandler(c echo.Context) error {
	var form UserLoginRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	var user User
	if err := DB.Preload("Schools").Preload("Tags").First(&user, &User{Email: form.Email}).Error; err != nil {
		// the user actually exists. quit register process
		LogService.Println("login: found user error: " + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusBadRequest, "no such user")
	}

	password, err := hex.DecodeString(user.Password)
	if err != nil {
		LogService.Println("login: failed to decode password bytes from hex string: " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if err := bcrypt.CompareHashAndPassword(password, []byte(form.Password)); err != nil {
		LogService.Println("login: password mismatch: " + err.Error() + spew.Sdump(user) + spew.Sdump(form))
		return echo.NewHTTPError(http.StatusUnauthorized, "password mismatch")
	}

	t, err := createJwt(&user)
	if err != nil {
		LogService.Println("login: failed to create jwt for user " + spew.Sdump(user) + ": " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, UserLoginResponse{
		Token:   t,
		Details: &user,
	})
}

func authVerifyEmailHandler(c echo.Context) error {
	key := c.Param("key")
	if key == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing verification key")
	}
	token, err := jwt.Parse(key, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "malformed jwt token supplied")
	}
	if token.Valid {
		return c.String(http.StatusOK, "successfully verified token")
	}
	return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
}

func authRegisterHandler(c echo.Context) error {
	var form UserRegisterRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	var existUser User
	if err := DB.First(&existUser, &User{Email: form.Email}).Error; err == nil {
		// the user actually exists. quit register process
		LogService.Println("register: user duplicated: existing " + spew.Sdump(existUser) + ", attempting " + spew.Sdump(form))
		return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		LogService.Println("register: failed to create hash from password: " + spew.Sdump(form))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create hash from password: use a longer password")
	}

	toSave := User{
		Email:    form.Email,
		Password: hex.EncodeToString(hashedPassword),

		FirstName: form.FirstName,
		LastName: form.LastName,
	}

	if form.Phone != nil {
		toSave.Phone = &sql.NullString{Valid: true, String: *form.Phone}
	}

	// calculate email verification token
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = toSave.Email
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["aud"] = "email-verify"
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		LogService.Println("register: failed to generate email verification token " + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate email verification token")
	}

	err = DB.Transaction(func(tx *gorm.DB) error {
		err := sendEmail(form.Email, "Teamer 账号注册邮箱验证", "RegisterEmailVerification", map[string]string{
			"User": form.Email,
			"Link": fmt.Sprintf("%s/api/auth/verify/email/%s", Conf.Server.Hostname, t),
		})
		if err != nil {
			LogService.Println("register: failed to send register email verification: " + spew.Sdump(err))
			//return echo.NewHTTPError(http.StatusInternalServerError, "failed to send verification email")
			return fmt.Errorf("failed to send verification email")
		}

		if DB.Create(&toSave).Error != nil {
			LogDb.Println("register: failed to create DB record: " + spew.Sdump(form))
			//return echo.NewHTTPError(http.StatusInternalServerError, "failed to create DB record")
			return fmt.Errorf("failed to create DB record")
		}
		return nil
	})

	if err != nil {
		LogService.Println("register: database transaction error " + spew.Sdump(toSave))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create user")
	}

	loginJwt, err := createJwt(&toSave)
	if err != nil {
		LogService.Println("register: failed to create jwt token for user " + spew.Sdump(toSave))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create jwt token")
	}

	return c.JSON(http.StatusAccepted, map[string]string{
		"token": loginJwt,
	})
}
