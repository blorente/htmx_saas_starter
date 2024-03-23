package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/middleware"
)

func getUserRecord(c echo.Context) (*models.Record, error) {
	info := apis.RequestInfo(c)
	userRecord := info.AuthRecord
	if userRecord == nil {
		return nil, fmt.Errorf("User not authenticated")
	}
	return userRecord, nil
}

// RegisterAuthRoutes registers the route group '/auth', which handles authentication.
func RegisterAuthRoutes(app *pocketbase.PocketBase, e *core.ServeEvent) {
	authGroup := e.Router.Group("/auth", middleware.LoadAuthContextFromCookie(app))

	authGroup.GET("/login", func(c echo.Context) error {
		app.Logger().Debug("GET Logging in")
		_, err := getUserRecord(c)
		if err == nil {
			app.Logger().Debug("User found. Redirecting")
			return c.Redirect(302, "/app/profile")
		}
		return c.File("public/components/login_form.html")
	})

	authGroup.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		app.Logger().Debug("Logging in: ", username)
		// TODO Actually log the user in, this just assumes it's in the DB
		token, err := lib.Login(e, username, password)
		if err != nil {
			app.Logger().Debug(fmt.Sprintf("Error logging in %s", err))
			c.Redirect(302, "/auth/login")
		}
		c.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    *token,
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
		})
		return c.Redirect(302, "/app/profile")
	})

	authGroup.POST("/logout", func(c echo.Context) error {
		app.Logger().Debug("Logging out")
		c.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    "",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			MaxAge:   -1,
		})
		return c.Redirect(302, "/login")
	})
}
