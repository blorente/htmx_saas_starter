package routes

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/middleware"
)

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
	authGroup.File("/login-form", "views/components/auth/login.html")

	authGroup.GET("/login", func(c echo.Context) error {
		_, err := lib.GetUserRecord(c)
		if err == nil {
			app.Logger().Debug("User found. Redirecting")
			return lib.NonHtmxRedirectToIndex(c)
		}
		return lib.RenderTemplate(c, registry, nil,
			"views/layout.html",
			"views/pages/login.html",
		)
	})

	authGroup.GET("/oauth-login/:provider", func(c echo.Context) error {
		provider := c.PathParams().Get("provider", "google")
		return lib.RenderTemplate(c, registry, AuthProviders[provider],
			"views/components/auth/login_with_provider.html",
		)
	})

	authGroup.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		app.Logger().Debug("Logging in: ", username)
		token, err := lib.LoginWithUsernameAndPassword(e, username, password)
		if err != nil {
			app.Logger().Debug(fmt.Sprintf("Error logging in %s", err))
			c.Redirect(302, "/auth/login")
		}
		lib.SetAuthCookie(c, *token)
		return lib.HtmxRedirectToIndex(c)
	})

	authGroup.POST("/logout", func(c echo.Context) error {
		app.Logger().Debug("Logging out")

		c.SetCookie(&http.Cookie{
			Name:     middleware.AuthCookieName,
			Value:    "",
			Path:     "/",
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})
		return lib.HtmxRedirectToIndex(c)
	})

	for _, formComponent := range []string{"username", "password"} {
		authGroup.GET(fmt.Sprintf("/%s", formComponent), func(c echo.Context) error {
			return lib.RenderTemplate(c, registry, nil,
				fmt.Sprintf("views/components/auth/%s.html", formComponent),
			)
		})
	}
}

// AuthRequestCallback will fire for every auth collection request.
// Here is where we'll set the cookie for oath authentication.
func AuthRequestCallback(e *core.RecordAuthEvent) error {
	lib.SetAuthCookie(e.HttpContext, e.Token)
	return nil
}
