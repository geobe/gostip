package controller

import (
	"errors"
	"github.com/geobe/gostip/go/model"
	"html"
	"net/http"
	"strconv"
	"time"
	"fmt"
	"log"
	"encoding/json"
)

// set an Applicants Data from http form parameters
func setApplicantData(app *model.Applicant, r *http.Request) {
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
	user, err := getUserFromSession(r)
	if err == nil {
		app.Data.Model.UpdatedBy = fmt.Sprintf("%s(%s)",user.Fullname, user.Login)
		app.Data.Model.Updater = user.ID
	} else {
		app.Data.Model.UpdatedBy = "self"
	}
}

func setEnrolledAt(app *model.Applicant) {
	if app.Data.EnrolledAt.IsZero() {
		app.Data.EnrolledAt = time.Now()
	}
}

func setResultData(app *model.Applicant, r *http.Request) {
	app.Data.Language = model.Lang(atoint(html.EscapeString(r.PostFormValue("language"))))
	val := html.EscapeString(r.PostFormValue("languageresult"))
	var f float32
	n, err := fmt.Sscanf(val, "%f", &f)
	if n == 1 && err == nil {
		app.Data.LanguageResult = int(f * 10.)
	} else {
		app.Data.LanguageResult = 0
	}
	for i := 0; i < model.NQESTION; i++ {
		rIndex := "r" + strconv.Itoa(i)
		val = html.EscapeString(r.PostFormValue(rIndex))
		n, err = fmt.Sscanf(val, "%f", &f)
		if n == 1 && err == nil {
			app.Data.Results[i] = int(f * 10.)
		} else {
			app.Data.Results[i] = 0
		}
	}
}

// fetch an applicant by its primary key that was submitted in form field "fieldname" and
// save applicant in session
func fetchApplicant(w http.ResponseWriter, r *http.Request, fieldname string, deleted ...bool) (app model.Applicant, err error) {
	if err = parseSubmission(w, r); err != nil {
		return
	}
	actionKey := html.EscapeString(r.PostFormValue("action"))
	appId, err := keyFromForm(w, r, fieldname)
	if err != nil {
		return
	}
	app, err = retrieveApplicant(appId, w, deleted...)
	if err != nil {
		return
	}
	if err = storeApplicant(w, r, app, actionKey); err != nil {
		return
	}
	return
}

// store current applicant the current session. the key parameter distinguishes
// between different applicant objects that are used simultaniously in the
// web interface, e.g. on different tab pages.
func storeApplicant(w http.ResponseWriter, r *http.Request, app model.Applicant, key string) (err error) {
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		log.Printf("SessionStore.get error %s", err)
		return
	}
	session.Values[S_APPLICANT + key] = app
	if err = session.Save(r, w); err != nil {
		log.Printf("SessionStore.save error %s", err)
		return
	}
	return
}

// parse form submission with error handling
func parseSubmission(w http.ResponseWriter, r *http.Request) (err error) {
	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("error %v, status %v\n", err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// get key from form with error handling
func keyFromForm(w http.ResponseWriter, r *http.Request, fieldname string) (id int, err error) {
	id, err = strconv.Atoi(html.EscapeString(r.PostFormValue(fieldname)))
	if err != nil {
		http.Error(w, "Conversion error: " + err.Error(), http.StatusInternalServerError)
		log.Printf("error %v, status %v\n", "Conversion error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

// retrieve an applicant from db with error handling
func retrieveApplicant(appId int, w http.ResponseWriter, deleted ...bool) (app model.Applicant, err error) {
	db := model.Db()
	if len(deleted) > 0 && deleted[0] {
		db = db.Unscoped()
	}
	db.Preload("Data").Preload("Data.Oblast").First(&app, appId)
	if app.ID == 0 {
		err = errors.New("Data integrity error")
		http.Error(w, "Data integrity error", http.StatusInternalServerError)
		log.Printf("error %v, status %v\n", "Data integrity error", http.StatusInternalServerError)
	}
	return
}

// get the previously stored applicant from the current session. the key parameter distinguishes
// between different applicant objects that are used simultaniously in the
// web interface, e.g. on different tab pages.
func applicantFromSession(key string, r *http.Request) (app model.Applicant, err error) {
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		return
	}
	val := session.Values[S_APPLICANT + key]
	var ok bool
	if app, ok = val.(model.Applicant); !ok {
		err = errors.New("conversion error from session")
		return
	}
	return
}

// save edited applicant data to database
func saveApplicantSubmission(w http.ResponseWriter, r *http.Request) {
	if err := parseSubmission(w, r); err != nil {
		http.Error(w, "Request parse error: " + err.Error(), http.StatusInternalServerError)
		log.Printf("error %v, status %v\n", "Request parse error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	action := html.EscapeString(r.PostFormValue("action"))
	app, err := applicantFromSession(action, r)
	if err != nil {
		http.Error(w, "Session store error: " + err.Error(), http.StatusInternalServerError)
		log.Printf("action %s, error %v, status %v\n", action, "Session store error: " + err.Error(), http.StatusInternalServerError)
		return
	}
	// save the base data in any case
	setApplicantData(&app, r)
	switch action {
	case "enrol":
		setEnrolledAt(&app)
	case "results":
		setResultData(&app, r)
	}
	if err := model.Db().Save(&app).Error; err != nil {
		var appModified model.Applicant
		var dataOld model.ApplicantData
		dataSubmitted := app.Data
		model.Db().Preload("Data").Preload("Data.Oblast").First(&appModified, app.ID)
		if appModified.ID == 0 {
			w.Header().Set("Content-Type", "application/json")
			json, _ := json.Marshal("Object was deleted")
			log.Printf("editing deleted Object, message is %s\n", string(json))
			w.Write(json)
			return
		}
		model.Db().Unscoped().First(&dataOld, dataSubmitted.ID)
		dataModified := appModified.Data
		merge, err := MergeDiff(&dataOld, &dataSubmitted, &dataModified, true, "form")
		if err != nil {
			http.Error(w, "Submission merge error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("error %v, status %v\n", "Submission merge error: " + err.Error(), http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(merge)
		if err != nil {
			log.Printf("Json marshalling error %v\n", err)
			http.Error(w, "Json marshalling error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("error %v, status %v\n", "Json marshalling error: " + err.Error(), http.StatusInternalServerError)
			return
		}
		if err := storeApplicant(w, r, appModified, action); err != nil {
			log.Printf("Session store error %v\n", err)
			http.Error(w, "Session store error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("error %v, status %v\n", "Session store error: " + err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

//// save edited test results
//func saveResultSubmission(w http.ResponseWriter, r *http.Request) {
//	app, err := fetchApplicant(w, r, "appid")
//	if err == nil {
//		setApplicantData(&app, r, false)
//		setResultData(&app, r)
//		model.Db().Save(&app)
//		w.WriteHeader(http.StatusNoContent)
//	}
//}

// allow only http POST methods
func checkMethodAllowed(method string, w http.ResponseWriter, r *http.Request) error {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("error %v, status %v\n", "Method not allowed", http.StatusMethodNotAllowed)
		return errors.New("Message not allowed: " + r.Method)
	}
	return nil
}

// set applicant data into viewmodel
func setViewModel(app model.Applicant, vmod viewmodel) {
	data := app.Data
	if app.UpdatedAt.Before(time.Unix(0.0, 0.0)) {
		vmod["updatedat"] = time.Now().UnixNano()
	} else {
		vmod["updatedat"] = app.UpdatedAt.UnixNano()
	}
	vmod["appid"] = app.ID
	vmod["lastname"] = data.LastName
	vmod["firstname"] = data.FirstName
	vmod["fathersname"] = data.FathersName
	vmod["lastnametx"] = data.LastNameTx
	vmod["firstnametx"] = data.FirstNameTx
	vmod["fathersnametx"] = data.FathersNameTx
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
