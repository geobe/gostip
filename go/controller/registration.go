package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"html/template"
	"net/http"
	"strconv"
)

func HandleRegistration(w http.ResponseWriter, r *http.Request) {
	db := model.Db()
	values := viewmodel{
		// HTMLAttr unescapes string for use as HTML Attribute
		"disabled":     template.HTMLAttr(""),
		"action":       "register",
		"oblasts":      model.Oblasts(),
		"displaystyle": "",
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
			OblastID:    oblastid(html.EscapeString(r.PostFormValue("district"))),
			OrtSum:      ortint(html.EscapeString(r.PostFormValue("ort"))),
			OrtMath:     ortint(html.EscapeString(r.PostFormValue("ortmath"))),
			OrtPhys:     ortint(html.EscapeString(r.PostFormValue("ortphys"))),
		}
		var xappdata model.ApplicantData
		db.First(&xappdata, "last_name = ? and first_name = ? "+
			"and fathers_name = ? and phone = ?",
			appdata.LastName, appdata.FirstName,
			appdata.FathersName, appdata.Phone)
		if xappdata.ID == 0 {
			app := model.Applicant{}
			app.Data = appdata
			db.Create(&app)
			values["disabled"] = template.HTMLAttr("disabled")
			values["displaystyle"] = "none"
		}
		view.Views().ExecuteTemplate(w, "registration", values)
	}

}

func ortint(s string) (v int16) {
	//println("ortint: '" + s + "'")
	i, err := strconv.Atoi(s)
	if err != nil {
		v = 0
	} else {
		v = int16(i)
	}
	return
}

func oblastid(s string) (id uint) {
	//println("oblastid: '" + s + "'")
	i, err := strconv.Atoi(s)
	if err != nil {
		id = 0
	} else {
		id = uint(i)
	}
	return
}
