package util

import (
	"html/template"
	"strings"
)

var TemplateFuncMap = template.FuncMap{
	"sub": func(i int, sub int) int {
		return i - sub
	},
	"trimgn": func(s string) string {
		return strings.ToLower(strings.ReplaceAll(s, " ", ""))
	},
}
