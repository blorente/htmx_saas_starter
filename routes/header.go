package routes

import (
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

// RegisterHeaderRoutes registers the route group '/header', which handles displaying data to the header.
func RegisterHeaderRoutes(app *pocketbase.PocketBase, e *core.ServeEvent, registry *template.Registry) {
	headerGroup := e.Router.Group("/header", middleware.LoadAuthContextFromCookie(app))

	headerGroup.GET("/loginstate", func(c echo.Context) error {
		userRecord := c.Get(apis.ContextAuthRecordKey)
		if userRecord == nil {
			return c.File("views/components/header/login.html")
		}
		app.Logger().Debug("Found user, displaying info")

		var user *models.Record = userRecord.(*models.Record)
		var props = map[string]any{}

		name := user.GetString("name")
		if name != "" {
			props["name"] = name
		} else {
			props["name"] = user.Username()
		}
		avatar, err := lib.GetFileUrl(user, "avatar")
		app.Logger().Debug("BL: Got avatar", avatar)
		if err == nil {
			// BL: THis should be an Error
			app.Logger().Debug("Failed to get url for avatar", err)
			props["avatar"] = avatar
		}

		app.Logger().Debug("BL: Rendering user info with props", props)
		html, err := registry.LoadFiles(
			"views/components/header/user_info.html",
		).Render(props)

		if err != nil {
			return apis.NewNotFoundError("", err)
		}

		return c.HTML(http.StatusOK, html)
	})
}
