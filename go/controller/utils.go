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
	"reflect"
	"github.com/justinas/nosurf"
	"regexp"
	"strings"
)

// set an Applicants Data from http form parameters
func setApplicantData(app *model.Applicant, r *http.Request) {
	setIfPosted(&app.Data.LastName, "lastname", r)
	setIfPosted(&app.Data.FirstName, "firstname", r)
	setIfPosted(&app.Data.FathersName, "fathersname", r)
	setIfPosted(&app.Data.Phone, "phone", r)
	setIfPosted(&app.Data.Email, "email", r)
	setIfPosted(&app.Data.Home, "home", r)
	setIfPosted(&app.Data.School, "school", r)
	setIfPosted(&app.Data.OrtSum, "ort", r)
	setIfPosted(&app.Data.OrtMath, "ortmath", r)
	setIfPosted(&app.Data.OrtPhys, "ortphys", r)
	var oblastID uint
	ok := setIfPosted(&oblastID, "district", r)
	if ok && app.Data.OblastID != oblastID {
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
		rIndex := "result" + strconv.Itoa(i)
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
func saveApplicantSubmission(w http.ResponseWriter, r *http.Request, maySetResults bool) {
	if err := parseSubmission(w, r); err != nil {
		http.Error(w, "Request parse error: " + err.Error(), http.StatusInternalServerError)
		log.Printf("error %v, status %v\n", "Request parse error: " + err.Error(),
			http.StatusInternalServerError)
		return
	}
	action := html.EscapeString(r.PostFormValue("action"))
	appsession, err := applicantFromSession(action, r)
	if err != nil {
		http.Error(w, "Session store error: " + err.Error(), http.StatusInternalServerError)
		log.Printf("action %s, error %v, status %v\n", action, "Session store error: " + err.Error(),
			http.StatusInternalServerError)
		return
	}
	// copy old applicant from session
	app := appsession
	// update data from form values
	setApplicantData(&app, r)
	switch action {
	case "enrol":
		setEnrolledAt(&app)
	case "results":
		// make sure that results may be modified by current user
		if maySetResults {
			setResultData(&app, r)
		}
	}
	if err := model.Db().Save(&app).Error; err != nil {
		var appModified model.Applicant
		dataOld := appsession.Data
		dataSubmitted := app.Data
		// read modified applicant from db
		model.Db().Preload("Data").Preload("Data.Oblast").First(&appModified, appsession.ID)
		if appModified.ID == 0 {
			w.Header().Set("Content-Type", "application/json")
			json, _ := json.Marshal("Object was deleted")
			log.Printf("editing deleted Object, message is %s\n", string(json))
			w.Write(json)
			return
		}
		dataModified := appModified.Data
		merge, err := MergeDiff(&dataOld, &dataSubmitted, &dataModified, true, "form")
		if err != nil {
			http.Error(w, "Submission merge error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("error %v, status %v\n", "Submission merge error: " + err.Error(),
				http.StatusInternalServerError)
			return
		}
		MergeScaleResults(merge, "result")
		json, err := json.Marshal(merge)
		if err != nil {
			log.Printf("Json marshalling error %v\n", err)
			http.Error(w, "Json marshalling error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("error %v, status %v\n", "Json marshalling error: " + err.Error(),
				http.StatusInternalServerError)
			return
		}
		// store in session
		if err := storeApplicant(w, r, appModified, action); err != nil {
			log.Printf("Session store error %v\n", err)
			http.Error(w, "Session store error: " + err.Error(), http.StatusInternalServerError)
			log.Printf("error %v, status %v\n", "Session store error: " + err.Error(),
				http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

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

// Only change values if there is a value in the PostForm
// We need some help from reflect to make it work
func setIfPosted(target interface{}, key string, r *http.Request) bool {
	_, ok := r.Form[key]
	if ok {
		val := html.EscapeString(r.PostFormValue(key))
		// settable references the variable target points to
		targetref := reflect.ValueOf(target).Elem()
		ok = targetref.CanSet()
		if ! ok {
			return false
		}
		switch target.(type) {
		case *string:
			targetref.SetString(val)
		case *uint:
			targetref.SetUint(uint64(atoint(val)))
		case *int16:
			targetref.SetInt(int64(atoint(val)))
		case *int:
			targetref.SetInt(int64(atoint(val)))
		default:
			log.Printf("cannot convert to type %v\n", target)
		}
	}
	return ok
}

// convert string with digits into int avoiding errors
func atoint(s string) (id int) {
	id, err := strconv.Atoi(s)
	if err != nil {
		id = 0
	}
	return
}

func checkForRegistration(r *http.Request) uint {
	var appId uint
	if v, ok := r.Form["appid"]; ok {
		appId = uint(atoint(html.EscapeString(v[0])))
	} else {
		appId = 0
	}
	if token, ok := r.Form["csrf_token"]; ok {
		if ! nosurf.VerifyToken(nosurf.Token(r), token[0]) {
			appId = 0
		}

	} else {
		appId = 0
	}
	return appId
}
var cleanSub = regexp.MustCompile("(-[A-Z][A-Z])|(;q=0.\\d)")
var splitter = regexp.MustCompile("(q=0.\\d),")
var quality =	regexp.MustCompile(".*q=0.([0-9])")
func PreferedLanguages(r *http.Request) []string {
	lnghdr := r.Header["Accept-Language"]
	slnghdr := strings.Split(splitter.ReplaceAllString(lnghdr[0], "$1|"), "|")
	var langs [4]string
	found := make(map[string]bool)
	i := 0
outer:	for _, v := range slnghdr {
		qlty, _ := strconv.Atoi(quality.ReplaceAllString(v, "$1"))
		if qlty < 5 {
			break
		}
		lgs := cleanSub.ReplaceAllString(v, "")
		for _, lg := range strings.Split(lgs, ",") {
			if i == 4 {
				break outer
			}
			if ! found[lg] {
				langs[i] = lg
				found[lg] = true
				i++
			}
		}
	}
	return  langs[0:i]
}
