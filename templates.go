package main

import (
	"embed"
	"html/template"
)

//go:embed templates
var templateFS embed.FS

var metricTpl = template.Must(template.ParseFS(
	templateFS,
	"templates/layout/base.html",
	"templates/admin/metrics.html",
))