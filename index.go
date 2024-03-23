package main

import (
	"log"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("public"), true))
		// e.Router.GET("/api/styles", apis.StaticDirectoryHandler(os.DirFS("public"), false))
		// e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("public"), false))
		e.Router.File("/components/sidebar", "public/components/layout/sidebar.html")
		e.Router.File("/components/header", "public/components/layout/header.html")
		e.Router.File("/components/dashboard", "public/components/layout/dashboard.html")
		e.Router.GET("/hello", func(c echo.Context) error {
			return c.HTML(http.StatusOK, "<h1>HsdrELLOsdfsdf </h1>")
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
