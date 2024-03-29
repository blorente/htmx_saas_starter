package lib

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/tools/template"
)

// RenderTemplate renders a template and writes it to the context.
// It's supposed to be "c.Render but with pocketbase's template tools"
func RenderTemplate(c echo.Context, registry *template.Registry, data any, templates ...string) error {
	html, err := registry.LoadFiles(
		templates...,
	).Render(data)

	if err != nil {
		return apis.NewNotFoundError("", err)
	}
	return c.HTML(http.StatusOK, html)
}
