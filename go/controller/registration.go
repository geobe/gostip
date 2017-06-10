package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
	"github.com/justinas/nosurf"
	"log"
)

// ShowRegistration is handler for registration form get requests:
// On GET, show an empty registration form, on POST save data into db
// and show a confirmation message.
func ShowRegistration(w http.ResponseWriter, r *http.Request) {
	db := model.Db()
	if err := parseSubmission(w, r); err != nil {
		return
	}
	i18nlanguage := view.PreferedLanguages(r) [0]
	values := viewmodel{
		// HTMLAttr unescapes string for use as HTML Attribute
		"disabled":     template.HTMLAttr(""),
		"action":       "register/submit",
		"oblasts":      view.OblastsI18n(i18nlanguage),
		"district":     0,
		"buttons":      true,
		"displaystyle": "block",
		"csrftoken":    nosurf.Token(r),
		"language": i18nlanguage,
	}
	if err := checkMethodAllowed(http.MethodGet, w, r); err == nil {
		appId := checkForRegistration(r)
		if appId > 0 {
			var app model.Applicant
			db.Preload("Data").Preload("Data.Oblast").First(&app, appId)
			if app.ID == appId {
				setViewModel(app, values)
			}
		}
		view.Views().ExecuteTemplate(w, "registration", values)
	}

}

// SubmitRegistration is handler for registration form post requests:
// On POST save data into db and show a confirmation message.
func SubmitRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		db := model.Db()
		if err := parseSubmission(w, r); err != nil {
			return
		}
		i18nlanguage := view.PreferedLanguages(r) [0]
		values := viewmodel{
			// HTMLAttr unescapes string for use as HTML Attribute
			"disabled":     template.HTMLAttr("disabled"),
			"action":       "register",
			"oblasts":      view.OblastsI18n(i18nlanguage),
			"buttons":      false,
			"thankyou":     true,
			"displaystyle": "none",
			"csrftoken":    nosurf.Token(r),
			"language":     i18nlanguage,
		}
		var app model.Applicant
		appId := atoint(html.EscapeString(r.PostFormValue("appid")))
		log.Printf("registration app id %d", appId)
		log.Print(*r)
		if appId > 0 {
			var err error
			app, err = retrieveApplicant(appId, w)
			// update registration
			if err != nil {
				return
			}
			setApplicantData(&app, r)
			db.Save(&app)
		} else {
			// new registration
			app = model.Applicant{}
			appdata := model.ApplicantData{}
			app.Data = appdata
			setApplicantData(&app, r)
			db.Create(&app)
		}
		setViewModel(app, values)
		view.Views().ExecuteTemplate(w, "registration", values)
	}

}
