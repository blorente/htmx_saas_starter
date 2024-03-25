package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"

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

type AuthProvider struct {
	Name        string
	DisplayName string
	LogoRoute   string
}

var AuthProviders = map[string]AuthProvider{
	"google": AuthProvider{
		Name:        "google",
		DisplayName: "Google",
		LogoRoute:   "/img/google_auth/round_dark.svg",
	},
}

// RegisterAuthRoutes registers the route group '/auth', which handles authentication.
func RegisterAuthRoutes(app *pocketbase.PocketBase, e *core.ServeEvent, registry *template.Registry) {
	authGroup := e.Router.Group("/auth", middleware.LoadAuthContextFromCookie(app))
	authGroup.GET("/login", func(c echo.Context) error {
		_, err := getUserRecord(c)
		if err == nil {
			app.Logger().Debug("User found. Redirecting")
			return c.Redirect(302, "/")
		}
		html, err := registry.LoadFiles(
			"views/layout.html",
			"views/pages/login.html",
		).Render(map[string]any{
			"needs_pocketbase": true,
		})
		if err != nil {
			app.Logger().Error(fmt.Sprintf("Error rendering template: %s", err))
			return apis.NewNotFoundError("Error rendering template", err)
		}
		return c.HTML(http.StatusOK, html)
	})

	authGroup.GET("/oauth-login/:provider", func(c echo.Context) error {
		provider := c.PathParams().Get("provider", "google")
		html, err := registry.LoadFiles(
			"views/components/oauth/login_with_provider.html",
		).Render(AuthProviders[provider])
		if err != nil {
			app.Logger().Error("Error rendering template", err)
			return apis.NewNotFoundError("", err)
		}
		return c.HTML(http.StatusOK, html)
	})

	authGroup.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		app.Logger().Debug("Logging in: ", username)
		// TODO Actually log the user in, this just assumes it's in the DB
		token, err := lib.LoginWithUsernameAndPassword(e, username, password)
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
		return c.Redirect(302, "/")
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
		return c.Redirect(302, "/")
	})
}

// AuthRequestCallback will fire for every auth collection request.
// Here is where we'll set the cookie for oath authentication.
func AuthRequestCallback(e *core.RecordAuthEvent) error {
	e.HttpContext.SetCookie(&http.Cookie{
		Name:     middleware.AuthCookieName,
		Value:    e.Token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}
