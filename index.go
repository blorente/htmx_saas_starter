package main

import (
	"log"
	"net/http"

	echo "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"

	"github.com/blorente/htmx_saas_starter/routes"
)

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		registry := template.NewRegistry()

		// Static HTML
		e.Router.Static("/*", "public")
		e.Router.File("/components/header", "views/components/layout/header.html")
		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(302, "/landing")
		})

		e.Router.GET("/landing", func(c echo.Context) error {
			html, err := registry.LoadFiles(
				"views/layout.html",
				"views/pages/landing.html",
			).Render(nil)

			if err != nil {
				return apis.NewNotFoundError("", err)
			}
			return c.HTML(http.StatusOK, html)
		})

		routes.RegisterAuthRoutes(app, e, registry)
		routes.RegisterAppRoutes(app, e)
		routes.RegisterHeaderRoutes(app, e, registry)
		return nil
	})

	app.OnRecordAuthRequest().Add(routes.AuthRequestCallback)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
