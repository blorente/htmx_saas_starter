package middleware

// This is copied wholesale from here:
//   https://github.com/gobeli/pocketbase-htmx/blob/main/middleware/auth.go

import (
	"fmt"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

const AuthCookieName = "Auth"

func LoadAuthContextFromCookie(app core.App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			app.Logger().Debug(fmt.Sprintf("BL: Reading from cookie with ctx %v", c))
			tokenCookie, err := c.Request().Cookie(AuthCookieName)
			app.Logger().Debug(fmt.Sprintf("BL: Token cookie is %v, err is %s", tokenCookie, err))
			if err != nil {
				return next(c)
			}

			token := tokenCookie.Value
			record, err := app.Dao().FindAuthRecordByToken(
				token,
				app.Settings().RecordAuthToken.Secret,
			)

			app.Logger().Debug(fmt.Sprintf("BL: Found record %v", record))
			if err != nil {
				return next(c)
			}

			app.Logger().Debug(fmt.Sprintf("BL: Setting the ContextAuthRecord to %v", record))
			c.Set(apis.ContextAuthRecordKey, record)
			return next(c)
		}
	}
}

func AuthGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		record := c.Get(apis.ContextAuthRecordKey)

		if record == nil {
			return c.Redirect(302, "/auth/login")
		}

		return next(c)
	}
}
