package controller

import (
	"errors"
	"github.com/geobe/gostip/go/model"
	"html"
	"net/http"
	"strconv"
	"time"
	"fmt"
)

// set an Applicants Data from http form parameters
func setApplicantData(app *model.Applicant, r *http.Request, enrol bool) {
	app.Data.LastName = html.EscapeString(r.PostFormValue("lastname"))
	app.Data.FirstName = html.EscapeString(r.PostFormValue("firstname"))
	app.Data.FathersName = html.EscapeString(r.PostFormValue("fathersname"))
	app.Data.Phone = html.EscapeString(r.PostFormValue("phone"))
	app.Data.Email = html.EscapeString(r.PostFormValue("email"))
	app.Data.Home = html.EscapeString(r.PostFormValue("home"))
	app.Data.School = html.EscapeString(r.PostFormValue("school"))
	app.Data.OrtSum = atoint16(html.EscapeString(r.PostFormValue("ort")))
	app.Data.OrtMath = atoint16(html.EscapeString(r.PostFormValue("ortmath")))
	app.Data.OrtPhys = atoint16(html.EscapeString(r.PostFormValue("ortphys")))
	oblastID := atouint(html.EscapeString(r.PostFormValue("district")))
	if app.Data.OblastID != oblastID {
		var o model.Oblast
		model.Db().First(&o, oblastID)
		app.Data.Oblast = o
		app.Data.OblastID = oblastID
	}
	if r.PostFormValue("schoolok") != "" {
		app.Data.SchoolOk = true
	}
	if r.PostFormValue("districtok") != "" {
		app.Data.OblastOk = true
	}
	if r.PostFormValue("ortok") != "" {
		app.Data.OrtOk = true
	}
	if enrol && app.Data.EnrolledAt.IsZero() {
		app.Data.EnrolledAt = time.Now()
	}
}

func setResultData(app *model.Applicant, r *http.Request) {
	app.Data.Language = model.Lang(atoint(html.EscapeString(r.PostFormValue("languages"))))
	app.Data.LanguageResult = atoint(html.EscapeString(r.PostFormValue("languageresult")))
	var f float32
	for i := 0; i < model.NQESTION; i++ {
		rIndex := "r" + strconv.Itoa(i)
		val := html.EscapeString(r.PostFormValue(rIndex))
		n, err := fmt.Sscanf(val, "%f", &f)
		if n == 1 && err == nil {
			app.Data.Results[i] = int(f * 10. )
		} else {
			app.Data.Results[i] = 0
		}
	}
}

// fetch an applicant by its primary key that was submitted in form field "fieldname".
func fetchApplicant(w http.ResponseWriter, r *http.Request, fieldname string) (app model.Applicant, err error) {
	if err = parseSubmission(w, r); err != nil {
		return
	}
	appId, err := keyFromForm(w, r, fieldname)
	if err != nil {
		return
	}
	app, err = retrieveApplicant(appId, w, r)
	return
}

// parse form submission with error handling
func parseSubmission(w http.ResponseWriter, r *http.Request) (err error) {
	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// get key from form with error handling
func keyFromForm(w http.ResponseWriter, r *http.Request, fieldname string) (id int, err error) {
	id, err = strconv.Atoi(html.EscapeString(r.PostFormValue(fieldname)))
	if err != nil {
		http.Error(w, "Conversion error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// retrieve an applicant from db with error handling
func retrieveApplicant(appId int, w http.ResponseWriter, r *http.Request) (app model.Applicant, err error) {
	db := model.Db()
	db.Preload("Data").Preload("Data.Oblast").First(&app, appId)
	if app.ID == 0 {
		err = errors.New("Data integrity error")
		http.Error(w, "Data integrity error", http.StatusInternalServerError)
	}
	return
}

// allow only http POST methods
func checkMethodAllowed(method string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return errors.New("Message not allowed: " + r.Method)
	}
	return nil
}

// set applicant data into viewmodel
func setViewModel(app model.Applicant, vmod viewmodel) {
	data := app.Data
	vmod["appid"] = app.ID
	vmod["lastname"] = data.LastName
	vmod["firstname"] = data.FirstName
	vmod["fathersname"] = data.FathersName
	vmod["phone"] = data.Phone
	vmod["email"] = data.Email
	vmod["home"] = data.Home
	vmod["school"] = data.School
	vmod["district"] = data.OblastID
	vmod["districtname"] = data.Oblast.Name
	vmod["ort"] = data.OrtSum
	vmod["ortmath"] = data.OrtMath
	vmod["ortphys"] = data.OrtPhys

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

// convert string with digits into int avoiding errors
func atoint(s string) (id int) {
	id, err := strconv.Atoi(s)
	if err != nil {
		id = 0
	}
	return
}
