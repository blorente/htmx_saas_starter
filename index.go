package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	echo "github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"

	"github.com/blorente/htmx_saas_starter/lib"
	"github.com/blorente/htmx_saas_starter/middleware"
)

func GetUserRecord(c echo.Context) (*models.Record, error) {
	info := apis.RequestInfo(c)
	userRecord := info.AuthRecord
	if userRecord == nil {
		return nil, fmt.Errorf("User not authenticated")
	}
	return userRecord, nil
}

func main() {
	app := pocketbase.New()

	// serves static files from the provided public dir (if exists)
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {

		authGroup := e.Router.Group("/auth", middleware.LoadAuthContextFromCookie(app))

		authGroup.GET("/login", func(c echo.Context) error {
			app.Logger().Debug("GET Logging in")
			_, err := GetUserRecord(c)
			if err == nil {
				app.Logger().Debug("User found. Redirecting")
				return c.Redirect(302, "/app/profile")
			}
			return c.File("public/components/login_form.html")
		})

		authGroup.POST("/login", func(c echo.Context) error {
			username := c.FormValue("username")
			password := c.FormValue("password")
			app.Logger().Debug("Logging in: ", username)
			// TODO Actually log the user in, this just assumes it's in the DB
			token, err := lib.Login(e, username, password)
			if err != nil {
				app.Logger().Debug(fmt.Sprintf("Error logging in %s", err))
				c.Redirect(302, "/auth/login")
			}
			c.SetCookie(&http.Cookie{
				Name:     middleware.AuthCookieName,
				Value:    *token,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
			})
			return c.Redirect(302, "/app/profile")
		})

		// Static HTML
		e.Router.GET("/*", apis.StaticDirectoryHandler(os.DirFS("public"), true))
		e.Router.File("/components/sidebar", "public/components/layout/sidebar.html")
		e.Router.File("/components/header", "public/components/layout/header.html")
		e.Router.File("/components/dashboard", "public/components/layout/dashboard.html")
		e.Router.GET("/hello", func(c echo.Context) error {
			return c.HTML(http.StatusOK, "<h1>HsdrELLOsdfsdf</h1>")
		})

		// App routes -> Have the auth guard
		appGroup := e.Router.Group("/app", middleware.LoadAuthContextFromCookie(app), middleware.AuthGuard)
		appGroup.GET("/profile", func(c echo.Context) error {
			var userRecord *models.Record = c.Get(apis.ContextAuthRecordKey).(*models.Record)
			return c.HTML(http.StatusOK, fmt.Sprintf(`<p>Welcome, %s<\p>`, userRecord.Username()))
		})
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
