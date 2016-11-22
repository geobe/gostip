// registration.go implements all handlers and helper functions
// for applicant registration.
package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
	"strconv"
)

// handler function for registration form get and post requests:
// On GET, show an empty registration form, on POST save data into db
// and show a confirmation message.
func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	db := model.Db()
	values := viewmodel{
		// HTMLAttr unescapes string for use as HTML Attribute
		"disabled":     template.HTMLAttr(""),
		"action":       "register",
		"oblasts":      model.Oblasts(),
		"district":     0,
		"buttons":      true,
		"displaystyle": "block",
	}
	if r.Method == http.MethodGet {
		view.Views().ExecuteTemplate(w, "registration", values)
	} else if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// do form reading field by field to avoid post attacks
		// with additional values,
		// escape to avoid malicious javascript code insertion
		appdata := model.ApplicantData{
			LastName:    html.EscapeString(r.PostFormValue("lastname")),
			FirstName:   html.EscapeString(r.PostFormValue("firstname")),
			FathersName: html.EscapeString(r.PostFormValue("fathersname")),
			Phone:       html.EscapeString(r.PostFormValue("phone")),
			Email:       html.EscapeString(r.PostFormValue("email")),
			Home:        html.EscapeString(r.PostFormValue("home")),
			School:      html.EscapeString(r.PostFormValue("school")),
			OblastID:    atouint(html.EscapeString(r.PostFormValue("district"))),
			OrtSum:      atoint16(html.EscapeString(r.PostFormValue("ort"))),
			OrtMath:     atoint16(html.EscapeString(r.PostFormValue("ortmath"))),
			OrtPhys:     atoint16(html.EscapeString(r.PostFormValue("ortphys"))),
		}
		var xappdata model.ApplicantData
		// already registered?
		db.First(&xappdata, "last_name = ? and first_name = ? "+"and phone = ?",
			appdata.LastName, appdata.FirstName, appdata.Phone)
		// if no, save and redisplay form
		if xappdata.ID == 0 {
			app := model.Applicant{}
			app.Data = appdata
			db.Create(&app)
			values["disabled"] = template.HTMLAttr("disabled")
			values["displaystyle"] = "none"
			values["thankyou"] = true
			values["buttons"] = false
			setViewModel(app, values)
			//values["district"] = appdata.OblastID
		}
		view.Views().ExecuteTemplate(w, "registration", values)
	}

}

// convert string with digits into int16
func atoint16(s string) (v int16) {
	i, err := strconv.Atoi(s)
	if err != nil {
		v = 0
	} else {
		v = int16(i)
	}
	return
}

// convert string with digits into uint
func atouint(s string) (id uint) {
	i, err := strconv.Atoi(s)
	if err != nil {
		id = 0
	} else {
		id = uint(i)
	}
	return
}
