package templates

import (
	"embed"
	"html/template"
	"strings"
)

//go:embed *.html
var EmbeddedTemplates embed.FS

// Template functions available in templates
var templateFuncs = template.FuncMap{
	"title": func(s string) string {
		if len(s) == 0 {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"trimSuffix": func(s, suffix string) string {
		return strings.TrimSuffix(s, suffix)
	},
}

// GetTemplates loads and returns all embedded HTML templates
func GetTemplates() (*template.Template, error) {
	tmpl := template.New("").Funcs(templateFuncs)
	return tmpl.ParseFS(EmbeddedTemplates, "*.html")
}