package controller

import (
	"net/http"
	"github.com/geobe/gostip/go/view"
	"github.com/geobe/gostip/go/model"
	"log"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	i18nlanguage := view.PreferedLanguages(r) [0]
	values := viewmodel {
		"language": i18nlanguage,
	}
	if err := checkMethodAllowed(http.MethodGet, w, r); err == nil {
		if err := parseSubmission(w, r); err != nil {
			return
		}
		appId := checkForRegistration(r)
		log.Printf("appId = %d", appId)
		if appId > 0 {
			var app model.Applicant
			app, err = retrieveApplicant(appId, w)
			// update registration
			if err != nil {
				return
			}
			setApplicantData(&app, r)
			setViewModel(app, values)
			values["districtname"] = view.I18n(app.Data.Oblast.Name, i18nlanguage)
			values["thankyou"] = "thx"
			makeMail("confirmation_mail", i18nlanguage, values)
		}
	}
	view.Views().ExecuteTemplate(w, "index", values)
}
