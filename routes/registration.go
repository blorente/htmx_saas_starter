package routes

import (
	"fmt"
	"net/http"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/middleware"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

// RegisterRegistraitionRoutes registers the route group '/registration', which handles showing and validating the registration form,
// as well as performing the actual registration.
func RegisterRegistrationRoutes(app *pocketbase.PocketBase, e *core.ServeEvent, registry *template.Registry) {
	group := e.Router.Group("/registration", middleware.LoadAuthContextFromCookie(app))

	group.GET("/register", func(c echo.Context) error {
		_, err := getUserRecord(c)
		if err == nil {
			app.Logger().Debug("User found. Redirecting")
			return c.Redirect(302, "/")
		}
		html, err := registry.LoadFiles(
			"views/layout.html",
			"views/pages/register.html",
		).Render(map[string]any{
			"needs_pocketbase": true,
		})
		if err != nil {
			app.Logger().Error(fmt.Sprintf("Error rendering template: %s", err))
			return apis.NewNotFoundError("Error rendering template", err)
		}
		return c.HTML(http.StatusOK, html)
	})

	// group.POST("/validate", func(c ec)

	group.POST("/register", func(c echo.Context) error {
		request := lib.NewRegisterUserRequestFromContext(c)
		err := request.Validate(app)
		if err != nil {
			return c.Redirect(302, "/registration/register")
		}
		_, err = lib.RegisterNewUser(app, &request)
		if err != nil {
			// TODO BL: Do better with error messages
			// Look at https://htmx.org/examples/inline-validation/
			c.Redirect(302, "/registration/register")
		}

		token, err := lib.LoginWithUsernameAndPassword(e, request.Username, request.Password)
		if err != nil {
			// TODO BL: Do better with error messages
			// Look at https://htmx.org/examples/inline-validation/

			app.Logger().Debug("BL: Error during registration, ", err)
			c.Redirect(302, "/registration/register")
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
}
