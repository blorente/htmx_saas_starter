package main

import (
	"log"
	"net/http"

	echo "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/blorente/htmx_saas_starter/routes"
)

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		// Static HTML
		e.Router.File("/components/sidebar", "public/components/layout/sidebar.html")
		e.Router.File("/components/header", "public/components/layout/header.html")
		e.Router.File("/components/dashboard", "public/components/layout/dashboard.html")
		e.Router.GET("/hello", func(c echo.Context) error {
			return c.HTML(http.StatusOK, "<h1>HsdrELLOsdfsdf</h1>")
		})

		routes.RegisterAuthRoutes(app, e)
		routes.RegisterAppRoutes(app, e)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
