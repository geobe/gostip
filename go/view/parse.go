// Parse all templates once and make them globally available
package view

import (
	"html/template"
	"os"
	"github.com/geobe/gostip/go/model"
	"strings"
)

const base = model.Base + "/go/view/*.go.html"

var views = Templates()

func Templates() *template.Template {
	pwd, _ := os.Getwd()
	pwd += "/"
	var path string
	if strings.HasSuffix(pwd, "main/") {
		path = pwd + "../view/*.go.html"
	} else {
		path = pwd + base
	}
	// Create a function map with dict and refer function
	funcs := map[string]interface{}{
		"dict": Dict,
		"adddict": AddDict,
		"mergedict": MergeDict,
		"safeatt": SafeAtt,
		"concat": Concat,
		"iskind": IsKind,
		"ismod": IsMod,
		"i18n": I18n,
	}

	// zuerst die function map in ein template parsen
	t := template.New("gostip.html").Funcs(funcs)
	// dann alle template files dazu parsen
	template.Must(t.ParseGlob(path))
	return t
}

func Views() *template.Template {
	return views
}
