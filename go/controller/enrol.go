// Package controller file enrol.go implements handler and helper functions for enrolment.
// I.e. for tabs enrol and edit.
package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
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
		return
	}
	values := viewmodel{
		"oblasts": model.Oblasts(),
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
		saveApplicantSubmission(w, r)
	}
}

// SubmitEdit is handler that accepts form submissions from the edit tab.
// Only http POST method is accepted.
func SubmitApplicantEdit(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		saveApplicantSubmission(w, r)
	}
}
