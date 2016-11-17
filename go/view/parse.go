// Parse all templates once and make them globally available
package view

import (
	"html/template"
	"os"
)

const base = "/github.com/geobe/gostip/go/view/*.html"

var views = Templates()

func Templates() *template.Template {
	pwd, _ := os.Getwd()
	t := template.Must(template.ParseGlob(pwd + base))
	return t
}

func Views() *template.Template {
	return views
}
