package controller

import (
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
	"github.com/justinas/nosurf"
	"log"
)

// ShowEnrol is handler to show the selected applicant from the
// search select element for enrol and edit tabs. It returns
// an html page fragment that is inserted into the respective tab area.
func ShowEnrol(w http.ResponseWriter, r *http.Request) {
	if checkMethodAllowed(http.MethodPost, w, r) != nil {
		return
	}
	action := html.EscapeString(r.PostFormValue("action"))
	app, err := fetchApplicant(w, r, "appid", false)
	if err != nil {
		log.Printf("fetch error %s", err)
		return
	}
	i18nlanguage := view.PreferedLanguages(r) [0]
	values := viewmodel{
		"oblasts": view.OblastsI18n(i18nlanguage),
		"csrftoken": nosurf.Token(r),
		"csrfid": "csrf_id_enrol",
		"language":     i18nlanguage,
	}
	setViewModel(app, values)

	if action == "enrol" {
		values["disabled"] = template.HTMLAttr("disabled='true'")
		view.Views().ExecuteTemplate(w, "work_enrol", values)
	} else {
		view.Views().ExecuteTemplate(w, "work_edit", values)
	}
}

// SubmitEnrol is handler that accepts form submissions from the enrol tab.
// Only http POST method is accepted.
func SubmitEnrol(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		saveApplicantSubmission(w, r, false)
	}
}

// SubmitEdit is handler that accepts form submissions from the edit tab.
// Only http POST method is accepted.
func SubmitApplicantEdit(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		saveApplicantSubmission(w, r, false)
	}
}
