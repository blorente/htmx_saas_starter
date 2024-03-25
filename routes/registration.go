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

	group.GET("/form", func(c echo.Context) error {
		html, err := renderFormTemplate(nil, registry)
		if err != nil {
			return err
		}
		return c.HTML(http.StatusOK, html)
	})

	group.POST("/register", func(c echo.Context) error {
		request := lib.NewRegisterUserRequestFromContext(c)
		err := request.Validate(app)
		if err != nil {
			html, err := renderFormTemplate(&err, registry)
			if err != nil {
				return err
			}
			return c.HTML(http.StatusOK, html)
		}
		_, err = lib.RegisterNewUser(app, &request)
		if err != nil {
			html, err := renderFormTemplate(&err, registry)
			if err != nil {
				return err
			}
			return c.HTML(http.StatusOK, html)
		}

		token, err := lib.LoginWithUsernameAndPassword(e, request.Username, request.Password)
		if err != nil {
			return c.Redirect(302, "/")
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

func renderFormTemplate(inErr *error, registry *template.Registry) (string, error) {
	props := make(map[string]any)

	if inErr != nil {
		props["error"] = fmt.Sprintf("%s", *inErr)
	}
	html, err := registry.LoadFiles(
		"views/components/registration/form.html",
	).Render(props)
	if err != nil {
		return "", apis.NewNotFoundError("Error rendering template", err)
	}
	return html, err
}
