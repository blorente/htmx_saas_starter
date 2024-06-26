package main

import (
	"log"

	echo "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/routes"
)

func main() {
	app := pocketbase.New()

	configPath := "config.yaml"

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		config, err := lib.NewConfigFromFile(configPath)
		log.Printf("Parsing config from %s", configPath)
		if err != nil {
			log.Fatal(err)
		}
		config.InitSettings(app)

		registry := template.NewRegistry()

		// Static HTML
		e.Router.Static("/*", "public")
		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(302, "/landing")
		})

		e.Router.GET("/landing", func(c echo.Context) error {
			return lib.RenderTemplate(c, registry, map[string]any{"Config": config},
				"views/layout.html",
				"views/pages/landing.html",
			)
		})

		routes.RegisterAuthRoutes(app, e, registry, *config)
		routes.RegisterRegistrationRoutes(app, e, registry, *config)
		routes.RegisterAppRoutes(app, e)
		routes.RegisterHeaderRoutes(app, e, registry)
		return nil
	})

	app.OnRecordAuthRequest().Add(routes.AuthRequestCallback)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
