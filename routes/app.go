package routes

import (
	"fmt"
	"net/http"

	"github.com/blorente/htmx_saas_starter/middleware"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

// RegisterAppRoutes registers the route group '/app', which is the entrypoint to all we can do once we're authenticated.
func RegisterAppRoutes(app *pocketbase.PocketBase, e *core.ServeEvent) {
	appGroup := e.Router.Group("/app", middleware.LoadAuthContextFromCookie(app), middleware.AuthGuard)
	appGroup.GET("/profile", func(c echo.Context) error {
		var userRecord *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
		return c.HTML(http.StatusOK, fmt.Sprintf(`<button class="button" hx-post="/auth/logout" hx-target="body">Welcome, %s</button>`, userRecord.Username()))
	})
}
