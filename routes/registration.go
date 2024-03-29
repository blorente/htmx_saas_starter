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
		_, err := lib.GetUserRecord(c)
		if err == nil {
			app.Logger().Debug("User found. Redirecting")
			return lib.HtmxRedirectToIndex(c)
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
			return lib.HtmxRedirectToIndex(c)
		}
		lib.SetAuthCookie(c, *token)
		return lib.HtmxRedirectToIndex(c)
	})

	group.POST("/validate-username", func(c echo.Context) error {
		app.Logger().Debug("BL: /validate-username")
		username := c.FormValue("username")
		err := lib.ValidateUsername(app, username)
		app.Logger().Debug("BL: Error is ", err)
		html, err := registry.LoadFiles(
			"views/components/registration/username.html",
		).Render(map[string]any{"error": err, "value": username})
		if err != nil {
			app.Logger().Error(fmt.Sprintf("Error rendering template: %s", err))
			return apis.NewNotFoundError("Error rendering template", err)
		}
		return c.HTML(http.StatusOK, html)
	})

	group.POST("/validate-email", func(c echo.Context) error {
		app.Logger().Debug(fmt.Sprintf("BL: /validate-email, ctx is %#v", c))
		email := c.FormValue("email")
		app.Logger().Debug(fmt.Sprintf("BL: /validate-email, email is %#v", email))
		err := lib.ValidateEmail(app, email)
		app.Logger().Debug(fmt.Sprintf("BL: Error is %s", err))
		html, err := registry.LoadFiles(
			"views/components/registration/email.html",
		).Render(map[string]any{"error": err, "value": email})
		if err != nil {
			app.Logger().Error(fmt.Sprintf("Error rendering template: %s", err))
			return apis.NewNotFoundError("Error rendering template", err)
		}
		return c.HTML(http.StatusOK, html)
	})

	group.POST("/validate-password", func(c echo.Context) error {
		app.Logger().Debug(fmt.Sprintf("BL: /validate-password, ctx is %#v", c))
		password := c.FormValue("password")
		repeatPassword := c.FormValue("repeat-password")
		err := lib.ValidatePassword(app, password, repeatPassword)
		app.Logger().Debug(fmt.Sprintf("BL: Error is %s", err))
		html, err := registry.LoadFiles(
			"views/components/registration/password.html",
		).Render(map[string]any{"error": err, "password": password, "repeat_password": repeatPassword})
		if err != nil {
			app.Logger().Error(fmt.Sprintf("Error rendering template: %s", err))
			return apis.NewNotFoundError("Error rendering template", err)
		}
		return c.HTML(http.StatusOK, html)
	})

	for _, formComponent := range []string{"email", "username", "password"} {
		group.GET(fmt.Sprintf("/%s", formComponent), func(c echo.Context) error {
			html, err := registry.LoadFiles(
				fmt.Sprintf("views/components/registration/%s.html", formComponent),
			).Render(nil)
			if err != nil {
				app.Logger().Error(fmt.Sprintf("Error rendering template: %s", err))
				return apis.NewNotFoundError("Error rendering template", err)
			}
			return c.HTML(http.StatusOK, html)
		})
	}

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
