package routes

import (
	"fmt"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/middleware"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
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
		return lib.RenderTemplate(c, registry,
			map[string]any{
				"needs_pocketbase": true,
			},
			"views/layout.html",
			"views/pages/register.html",
		)
	})

	group.GET("/form", func(c echo.Context) error {
		return renderRegistrationForm(c, registry, nil)
	})

	group.POST("/register", func(c echo.Context) error {
		request := lib.NewRegisterUserRequestFromContext(c)
		err := request.Validate(app)
		if err != nil {
			return renderRegistrationForm(c, registry, &err)
		}
		_, err = lib.RegisterNewUser(app, &request)
		if err != nil {
			return renderRegistrationForm(c, registry, &err)
		}

		token, err := lib.LoginWithUsernameAndPassword(e, request.Username, request.Password)
		if err != nil {
			return lib.HtmxRedirectToIndex(c)
		}
		lib.SetAuthCookie(c, *token)
		return lib.HtmxRedirectToIndex(c)
	})

	group.POST("/validate-username", func(c echo.Context) error {
		username := c.FormValue("username")
		err := lib.ValidateUsername(app, username)
		return lib.RenderTemplate(c, registry,
			map[string]any{"error": err, "value": username},
			"views/components/registration/username.html",
		)
	})

	group.POST("/validate-email", func(c echo.Context) error {
		email := c.FormValue("email")
		err := lib.ValidateEmail(app, email)
		return lib.RenderTemplate(c, registry,
			map[string]any{"error": err, "value": email},
			"views/components/registration/email.html",
		)
	})

	group.POST("/validate-password", func(c echo.Context) error {
		password := c.FormValue("password")
		repeatPassword := c.FormValue("repeat-password")
		err := lib.ValidatePassword(app, password, repeatPassword)
		return lib.RenderTemplate(c, registry,
			map[string]any{"error": err, "password": password, "repeat_password": repeatPassword},
			"views/components/registration/password.html",
		)
	})

	for _, formComponent := range []string{"email", "username", "password"} {
		group.GET(fmt.Sprintf("/%s", formComponent), func(c echo.Context) error {
			return lib.RenderTemplate(c, registry, nil,
				fmt.Sprintf("views/components/registration/%s.html", formComponent),
			)
		})
	}
}

func renderRegistrationForm(c echo.Context, registry *template.Registry, inErr *error) error {
	props := make(map[string]any)

	if inErr != nil {
		props["error"] = fmt.Sprintf("%s", *inErr)
	}
	return lib.RenderTemplate(c, registry, props,
		"views/components/registration/form.html",
	)
}
