package routes

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/template"

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

		var user *models.Record = userRecord.(*models.Record)
		name := user.Username()

		html, err := registry.LoadFiles(
			"views/components/header/user_info.html",
		).Render(map[string]any{
			"name": name,
		})
		if err != nil {
			// or redirect to a dedicated 404 HTML page
			return apis.NewNotFoundError("", err)
		}

		return c.HTML(http.StatusOK, html)
	})
}
