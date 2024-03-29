package lib

// Adapted from https://github.com/gobeli/pocketbase-htmx/blob/main/lib/auth.go

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/blorente/htmx_saas_starter/middleware"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tokens"
)

func LoginWithUsernameAndPassword(e *core.ServeEvent, username string, password string) (*string, error) {
	user, err := e.App.Dao().FindAuthRecordByUsername("users", username)
	if err != nil {
		return nil, fmt.Errorf("Login failed")
	}

	valid := user.ValidatePassword(password)
	if !valid {
		return nil, fmt.Errorf("Login failed")
	}

	s, tokenErr := tokens.NewRecordAuthToken(e.App, user)
	if tokenErr != nil {
		return nil, fmt.Errorf("Login failed")
	}

	return &s, nil
}

type RegisterNewUserRequest struct {
	Username       string
	Name           string
	Email          string
	Password       string
	RetypePassword string
}

func NewRegisterUserRequestFromContext(c echo.Context) RegisterNewUserRequest {
	return RegisterNewUserRequest{
		Username:       c.FormValue("username"),
		Name:           c.FormValue("name"),
		Email:          c.FormValue("email"),
		Password:       c.FormValue("password"),
		RetypePassword: c.FormValue("repeat-password"),
	}
}

func (r *RegisterNewUserRequest) Validate(app *pocketbase.PocketBase) error {
	var err error

	app.Logger().Debug(fmt.Sprintf("BL: Validating request %v", r))

	if err = ValidatePassword(app, r.Password, r.RetypePassword); err != nil {
		return err
	}

	if err = ValidateEmail(app, r.Email); err != nil {
		return err
	}

	if err = ValidateUsername(app, r.Username); err != nil {
		return err
	}
	return nil
}

func ValidatePassword(app *pocketbase.PocketBase, password string, retype string) error {
	if password != retype {
		return fmt.Errorf("Passwords do not match")
	}
	return nil
}

func ValidateUsername(app *pocketbase.PocketBase, username string) error {
	_, err := app.Dao().FindFirstRecordByData("users", "username", username)
	if err == nil {
		return fmt.Errorf("Username is already taken")
	}
	return nil
}

func ValidateEmail(app *pocketbase.PocketBase, email string) error {
	parsedEmail, err := mail.ParseAddress(email)
	fmt.Println("BL: Parsed email is %#v, err is %s", parsedEmail, err)
	if err != nil || parsedEmail.Address != email {
		fmt.Println("BL: err is %s", parsedEmail, err)
		return fmt.Errorf("Email is not a valid email address")
	}
	_, err = app.Dao().FindFirstRecordByData("users", "email", email)
	if err == nil {
		return fmt.Errorf("Email is already taken")
	}
	return nil
}

func RegisterNewUser(app *pocketbase.PocketBase, req *RegisterNewUserRequest) (*models.Record, error) {
	users, err := app.Dao().FindCollectionByNameOrId("users")
	if err != nil {
		return nil, err
	}

	newUserRecord := models.NewRecord(users)
	err = newUserRecord.SetUsername(req.Username)
	if err != nil {
		return nil, err
	}

	err = newUserRecord.SetEmail(req.Email)
	if err != nil {
		return nil, err
	}

	err = newUserRecord.SetPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		newUserRecord.Set("name", req.Name)
	} else {
		newUserRecord.Set("name", req.Username)
	}

	err = newUserRecord.RefreshTokenKey()
	if err != nil {
		return nil, err
	}

	err = app.Dao().SaveRecord(newUserRecord)
	if err != nil {
		return nil, err
	}
	return newUserRecord, nil
}

func GetUserRecord(c echo.Context) (*models.Record, error) {
	info := apis.RequestInfo(c)
	userRecord := info.AuthRecord
	if userRecord == nil {
		return nil, fmt.Errorf("User not authenticated")
	}
	return userRecord, nil
}

func SetAuthCookie(c echo.Context, token string) {
	c.SetCookie(&http.Cookie{
		Name:     middleware.AuthCookieName,
		Value:    token,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
	})
}
