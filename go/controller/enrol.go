// enrol.go implements handler and helper functions for enrolment.
// I.e. for tabs enrol and edit.
package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"github.com/pkg/errors"
	"html"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

// Handler function to show the selected applicant from the
// search select element for enrol and edit tabs. It returns
// an html page fragment that is inserted into the respective tab area.
func ShowEnrol(w http.ResponseWriter, r *http.Request) {
	if err := onlyPostAllowed(w, r); err != nil {
		return
	}
	app, err := fetchApplicant(w, r, "appid")
	if err != nil {
		return
	}
	action := html.EscapeString(r.PostFormValue("action"))
	enrol := action == "enrol"

	data := app.Data
	values := viewmodel{
		"enrolaction":  "/enrol/submit",
		"displaystyle": "none",
		"applicantid":  app.ID,
		"lastname":     data.LastName,
		"firstname":    data.FirstName,
		"fathersname":  data.FathersName,
		"phone":        data.Phone,
		"email":        data.Email,
		"home":         data.Home,
		"school":       data.School,
		"oblasts":      model.Oblasts(),
		"district":     data.Oblast.ID,
		"districtname": data.Oblast.Name,
		"ort":          data.OrtSum,
		"ortmath":      data.OrtMath,
		"ortphys":      data.OrtPhys,
	}
	//if err = addRoles(r, values); err != nil {
	//	return
	//}
	if enrol {
		values["disabled"] = template.HTMLAttr("disabled='true'")
		view.Views().ExecuteTemplate(w, "work_enrol", values)
	} else {
		view.Views().ExecuteTemplate(w, "work_edit", values)
	}
}

// Handler function that accepts form submissions from the enrol tab.
// Only http POST method is accepted.
func SubmitEnrol(w http.ResponseWriter, r *http.Request) {
	saveSubmission(w, r, true)
}

// Handler function that accepts form submissions from the edit tab.
// Only http POST method is accepted.
func SubmitEdit(w http.ResponseWriter, r *http.Request) {
	saveSubmission(w, r, false)
}

// save edited applicant data to database
func saveSubmission(w http.ResponseWriter, r *http.Request, enrol bool) {
	if err := onlyPostAllowed(w, r); err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	app, err := fetchApplicant(w, r, "applicantid")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// do form reading field by field to avoid post attacks
	// with additional values,
	// escape to avoid malicious javascript code insertion
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
	if enrol {
		app.Data.SchoolOk = r.PostFormValue("schoolok") != ""
		app.Data.OblastOk = r.PostFormValue("district") != ""
		app.Data.OrtOk = r.PostFormValue("ortok") != ""
	}

	oblastID := atouint(html.EscapeString(r.PostFormValue("district")))
	if app.Data.EnrolledAt.IsZero() {
		app.Data.EnrolledAt = time.Now()
	}
	if app.Data.OblastID != oblastID {
		var o model.Oblast
		model.Db().First(&o, oblastID)
		app.Data.Oblast = o
	}
	model.Db().Save(&app)
	w.WriteHeader(http.StatusNoContent)

}

// allow only http POST methods
func onlyPostAllowed(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return errors.New("Message not allowed: " + r.Method)
	}
	return nil
}

// fetch an applicant by its primary key that was submitted in form field "fieldname".
func fetchApplicant(w http.ResponseWriter, r *http.Request, fieldname string) (app model.Applicant, err error) {
	if err = r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	appId, err := strconv.Atoi(html.EscapeString(r.PostFormValue(fieldname)))
	if err != nil {
		http.Error(w, "Conversion error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	db := model.Db()
	db.Preload("Data").Preload("Data.Oblast").First(&app, appId)
	if app.ID == 0 {
		err = errors.New("Data integrity error")
		http.Error(w, "Data integrity error", http.StatusInternalServerError)
	}
	return
}
