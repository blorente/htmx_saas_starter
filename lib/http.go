package lib

import "github.com/labstack/echo/v5"

func NonHtmxRedirectToIndex(c echo.Context) error {
	return c.Redirect(302, "/")
}
func HtmxRedirectToIndex(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/")
	c.Response().Header().Set("HX-Push-Url", "true")
	c.Response().Header().Set("HX-Refresh", "true")
	return c.Redirect(200, "/")
}

func NoCacheResponse(c echo.Context) {
	c.Response().Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Response().Header().Set("Pragma", "no-cache")
	c.Response().Header().Set("Expires", "0")
}
