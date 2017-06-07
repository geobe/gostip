package main

import (
	"log"
	"os"
	"text/template"
	"github.com/geobe/gostip/go/view"
	"testing"
)

type District struct {
	ID int
	Name string
}

// this test demonstrates how dot access can be passed into an inner template
func TestReferencePassing(t *testing.T) {
	// Define a template.
	const tpl1 = `
{{define "outer"}}

look at dot (a slice or array) {{.}} and some indexed value {{index . 3}}

Reference to array with index (must be passed as a string!) {{refer . "1"}}
Reference to named field of a struct {{refer (index . 3) "Name"}}

	{{/*now expand an inner template and pass it a newly constructed map (with function dict) */}}
	{{template "inner" dict "dot" . "animal" (index . 1) "where" ((index . 3).Name) "key" "animal" -}}
{{end}}

{{define "inner" -}}
	It'a an {{.animal -}}! Where is it? {{.where}}
	outer dot is here, part o a map: {{.dot}}

	Here is what reference function delivers: {{(refer . "key")}}
	You can pass a value into a nested template that selects which value of "Dot" is to be used
	Reference to map {{refer . (.key)}}
	With built-in function index {{index . (.key)}}
{{end -}}
`
	// Create a function map with dict and refer function
	funcs := map[string]interface{}{
		"dict": view.Dict,
		//"refer": view.DotReference,
	}
	// Create a new template and add function map.
	tpl := template.New("test").Funcs(funcs)
	// parse template
	template.Must(tpl.Parse(tpl1))
	s := []interface{}{"hi", "Elk", "in", District{Name: "Saxony", ID: 42}}

	err := tpl.ExecuteTemplate(os.Stdout, "outer", s)
	if err != nil {
		log.Println("executing template:", err)
	}
}

