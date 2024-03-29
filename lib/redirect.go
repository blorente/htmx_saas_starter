package lib

import "github.com/labstack/echo/v5"

func NonHtmxRedirectToIndex(c echo.Context) error {
	return c.Redirect(302, "/")
}
func HtmxRedirectToIndex(c echo.Context) error {
	c.Response().Header().Set("HX-Redirect", "/")
	return c.Redirect(200, "/")
}
