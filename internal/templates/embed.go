package templates

import (
	"embed"
	"html/template"
)

//go:embed layout/*.html admin/*.html
var fs embed.FS

var MetricTpl = template.Must(template.ParseFS(
	fs,
	"layout/base.html",
	"admin/metrics.html",
))