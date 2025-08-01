package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html
var EmbeddedTemplates embed.FS

// GetTemplates loads and returns all embedded HTML templates
func GetTemplates() (*template.Template, error) {
	return template.ParseFS(EmbeddedTemplates, "*.html")
}