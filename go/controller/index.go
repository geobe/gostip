package controller

import (
	"net/http"
	"github.com/geobe/gostip/go/view"
	"github.com/geobe/gostip/go/model"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	values := viewmodel {}
	if err := checkMethodAllowed(http.MethodGet, w, r); err == nil {
		if err := parseSubmission(w, r); err != nil {
			return
		}
		appId := checkForRegistration(r)
		if appId > 0 {
			db := model.Db()
			var app model.Applicant
			db.Preload("Data").First(&app, appId)
			if app.ID == appId {
				values["thankyou"] = "thx"
				values["lastname"] = app.Data.LastName
				values["firstname"]  = app.Data.FirstName
				values["email"] = app.Data.Email
			}
		}
	}
	values["language"] = view.PreferedLanguages(r) [0]
	view.Views().ExecuteTemplate(w, "index", values)
}
