package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/middleware"
)

// RegisterHeaderRoutes registers the route group '/header', which handles displaying data to the header.
func RegisterHeaderRoutes(app *pocketbase.PocketBase, e *core.ServeEvent, registry *template.Registry) {
	headerGroup := e.Router.Group("/header", middleware.LoadAuthContextFromCookie(app))

	headerGroup.GET("/loginstate", func(c echo.Context) error {
		user, err := lib.GetUserRecord(c)
		if err != nil {
			return c.File("views/components/header/login.html")
		}
		app.Logger().Debug("Found user, displaying info")

		var props = map[string]any{}

		name := user.GetString("name")
		if name != "" {
			props["name"] = name
		} else {
			props["name"] = user.Username()
		}
		avatar, err := lib.GetFileUrl(user, "avatar")
		if err == nil {
			props["avatar"] = avatar
		}

		html, err := registry.LoadFiles(
			"views/components/header/user_info.html",
		).Render(props)

		if err != nil {
			return apis.NewNotFoundError("", err)
		}

		return c.HTML(http.StatusOK, html)
	})
}
