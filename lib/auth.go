package lib

// Copied wholesale from https://github.com/gobeli/pocketbase-htmx/blob/main/lib/auth.go

import (
	"fmt"
	"net/mail"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
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

	app.Logger().Debug(fmt.Sprintf("BL: Validating request %v", r))
	if r.Password != r.RetypePassword {
		return fmt.Errorf("Passwords do not match")
	}
	parsedEmail, err := mail.ParseAddress(r.Email)
	if err != nil || parsedEmail.Address != r.Email {
		return fmt.Errorf("Email is not a valid email address")
	}
	_, err = app.Dao().FindFirstRecordByData("users", "email", r.Email)
	if err == nil {
		return fmt.Errorf("Email is already taken")
	}
	_, err = app.Dao().FindFirstRecordByData("users", "username", r.Username)
	if err == nil {
		return fmt.Errorf("Username is already taken")
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
