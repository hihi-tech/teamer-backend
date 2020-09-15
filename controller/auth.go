package controller

import (
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"teamer/model"
	"time"
)

type UserLoginRequestForm struct {
	Email string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserAuthResponse struct {
	Token string       `json:"token"`
	Details *model.User `json:"details"`
}

type UserRegisterRequestForm struct {
	Email string `json:"email" validate:"email,max=64,required"`
	Password string `json:"password" validate:"min=8,max=64,required"`
	Phone    *string `json:"phone" validate:"max=18,required"`

	FirstName string `json:"firstName" validate:"max=64,required"`
	LastName  string `json:"lastName" validate:"max=64,required"`
	Birthday  string `json:"birthday" validate:"max=24,required"`

	Schools []uint `json:"schools"`
	//Tags []uint `json:"tags"`
}

func (ct Controller) CreateJwt(user *model.User) (string, error) {
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

// Login godoc
// @Summary Login
// @Description Login using Teamer credentials
// @Tags accounts
// @Accept json
// @Produce  json
// @Success 200 {object} UserAuthResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/login [post]
func (ct Controller) Login(c echo.Context) error {
	var form UserLoginRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	var user model.User
	if err := ct.db.Preload("Schools").Preload("Tags").First(&user, &model.User{Email: form.Email}).Error; err != nil {
		// the user actually exists. quit register process
		ct.logger.Println("login: found user error: " + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusBadRequest, "no such user")
	}

	password, err := hex.DecodeString(user.Password)
	if err != nil {
		ct.logger.Println("login: failed to decode password bytes from hex string: " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	if err := bcrypt.CompareHashAndPassword(password, []byte(form.Password)); err != nil {
		ct.logger.Println("login: password mismatch: " + err.Error() + spew.Sdump(user) + spew.Sdump(form))
		return echo.NewHTTPError(http.StatusUnauthorized, "电子邮箱或密码错误")
	}

	t, err := ct.CreateJwt(&user)
	if err != nil {
		ct.logger.Println("login: failed to create jwt for user " + spew.Sdump(user) + ": " + err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "internal error")
	}

	return c.JSON(http.StatusOK, UserAuthResponse{
		Token:   t,
		Details: &user,
	})
}

// VerifyEmail godoc
// @Summary Verify Email
// @Description The landing page after user clicks the verify button in the register email
// @Tags accounts
// @Produce json
// @Param key path string true "Verify Key"
// @Success 200 {string}
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/verify/email/{key} [get]
func (ct Controller) VerifyEmail(c echo.Context) error {
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

// Register godoc
// @Summary Register
// @Description Register by using details provided
// @Tags accounts
// @Accept json
// @Produce json
// @Success 200 {object} UserAuthResponse
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/register [post]
func (ct Controller) Register(c echo.Context) error {
	var form UserRegisterRequestForm
	if err := c.Bind(&form); err != nil {
		return DefaultBadRequestResponse
	}
	if err := c.Validate(&form); err != nil {
		return DefaultBadRequestResponse
	}

	var existUser model.User
	if err := ct.db.First(&existUser, &model.User{Email: form.Email}).Error; err == nil {
		// the user actually exists. quit register process
		ct.logger.Println("register: user duplicated: existing " + spew.Sdump(existUser) + ", attempting " + spew.Sdump(form))
		return echo.NewHTTPError(http.StatusBadRequest, "user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		ct.logger.Println("register: failed to create hash from password: " + spew.Sdump(form))
		return echo.NewHTTPError(http.StatusBadRequest, "failed to create hash from password: use a longer password")
	}

	var schools []*model.School
	for _, school := range form.Schools {
		var foundSchool model.School
		if err := ct.db.First(&foundSchool, school).Error; err != nil {
			spew.Dump(err)
			return echo.NewHTTPError(http.StatusBadRequest, "cannot found school with id " + strconv.Itoa(int(school)))
		}
		schools = append(schools, &foundSchool)
	}

	toSave := model.User{
		Email:    form.Email,
		Password: hex.EncodeToString(hashedPassword),
		Birthday: form.Birthday,
		FirstName: form.FirstName,
		LastName: form.LastName,
		Phone: form.Phone,
		Schools: schools,
	}

	spew.Dump(form, toSave)

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
		ct.logger.Println("register: failed to generate email verification token " + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate email verification token")
	}

	err = ct.db.Transaction(func(tx *gorm.DB) error {
		err := ct.SendEmail(form.Email, "Teamer 账号注册邮箱验证", "RegisterEmailVerification", map[string]string{
			"User": form.Email,
			"Link": fmt.Sprintf("%s/api/auth/verify/email/%s", ct.config.Server.Hostname, t),
		})
		if err != nil {
			ct.logger.Println("register: failed to send register email verification: " + spew.Sdump(err))
			//return echo.NewHTTPError(http.StatusInternalServerError, "failed to send verification email")
			return fmt.Errorf("failed to send verification email")
		}

		if err := ct.db.Create(&toSave).Error; err != nil {
			ct.logger.Println("register: failed to create db record: " + spew.Sdump(err))
			//return echo.NewHTTPError(http.StatusInternalServerError, "failed to create db record")
			return fmt.Errorf("failed to create db record")
		}
		return nil
	})

	if err != nil {
		ct.logger.Println("register: database transaction error " + spew.Sdump(toSave))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user")
	}

	loginJwt, err := ct.CreateJwt(&toSave)
	if err != nil {
		ct.logger.Println("register: failed to create jwt token for user " + spew.Sdump(toSave))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create jwt token")
	}

	var savedUser model.User
	if err := ct.db.First(&savedUser, toSave.ID).Error; err != nil {
		ct.logger.Println("register: failed to get saved user " + spew.Sdump(toSave) + spew.Sdump(err))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get user details")
	}

	return c.JSON(http.StatusAccepted, UserAuthResponse{
		Token:   loginJwt,
		Details: &savedUser,
	})
}
