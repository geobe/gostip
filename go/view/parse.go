// Parse all templates once and make them globally available
package view

import (
	"html/template"
	"os"
	"github.com/geobe/gostip/go/model"
	"strings"
)

const base = model.Base + "/go/view/*.html"

var views = Templates()

func Templates() *template.Template {
	pwd, _ := os.Getwd()
	pwd += "/"
	var path string
	if strings.HasSuffix(pwd, "main/") {
		path = pwd + "../view/*.html"
	} else {
		path = pwd + base
	}
	//t := template.Must(template.ParseGlob(pwd + base))
	// Create a function map with dict and refer function
	funcs := map[string]interface{}{
		"dict": Dict,
		"adddict": AddDict,
		"mergedict": MergeDict,
		//"refer": DotReference,
		"safeatt": SafeAtt,
		"concat": Concat,
	}

	t := template.New("gostip.html").Funcs(funcs)
	template.Must(t.ParseGlob(path))
	return t
}

func Views() *template.Template {
	return views
}
