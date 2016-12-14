// Package controller file enrol.go implements handler and helper functions for enrolment.
// I.e. for tabs enrol and edit.
package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
	"fmt"
)

// ShowEnrol is handler to show the selected applicant from the
// search select element for enrol and edit tabs. It returns
// an html page fragment that is inserted into the respective tab area.
func ShowEnrol(w http.ResponseWriter, r *http.Request) {
	if checkMethodAllowed(http.MethodPost, w, r) != nil {
		return
	}
	app, err := fetchApplicant(w, r, "appid")
	if err != nil {
		return
	}
	action := html.EscapeString(r.PostFormValue("action"))
	enrol := action == "enrol"

	//data := app.Data
	values := viewmodel{
		//"enrolaction":  "/enrol/submit",
		//"displaystyle": "none",
		"oblasts": model.Oblasts(),
	}
	setViewModel(app, values)
	if enrol {
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
		saveApplicantSubmission(w, r, true)
	}
}

// SubmitEdit is handler that accepts form submissions from the edit tab.
// Only http POST method is accepted.
func SubmitApplicantEdit(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		saveApplicantSubmission(w, r, false)
	}
}

// save edited applicant data to database
func saveApplicantSubmission(w http.ResponseWriter, r *http.Request, enrol bool) {
	app, err := fetchApplicant(w, r, "appid")
	upat, err2 := keyFromForm(w, r, "updatedat")
	updb := app.UpdatedAt.UnixNano()
	fmt.Printf("updated at from form: %d, from db: %d, delta: %d\n", upat, updb, int64(upat) - updb)
	if err == nil && err2 == nil {
		setApplicantData(&app, r, enrol)
		if err := model.Db().Save(&app).Error; err != nil {
			var appModified model.Applicant
			var dataOld model.ApplicantData
			dataSubmitted := app.Data
			model.Db().Preload("Data").Preload("Data.Oblast").First(&appModified, app.ID)
			model.Db().Unscoped().First(&dataOld, dataSubmitted.ID)
			dataModified := appModified.Data

			merge, err := MergeDiff(&dataOld, &dataSubmitted, &dataModified, true, "form")
			if err != nil {
				fmt.Printf("error in merging %v\n", err)
			} else {
				for k, v := range merge {
					cnf := "update"
					ic := ""
					if v.Conflict {
						cnf = "conflict"
						ic = ">"
					}
					fmt.Printf("merge auto [%s]: %v <-%s %v is %s\n", k, v.Mine, ic, v.Other, cnf)
				}
			}
		} else {

		}
		fmt.Fprint(w, "hurz")
	}
}
