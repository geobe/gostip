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

func ShowEnrol(w http.ResponseWriter, r *http.Request) {
	if err := onlyPostAllowed(w, r); err != nil {
		return
	}
	app, err := fetchApplicant(w, r, "appid")
	if err != nil {
		return
	}
	data := app.Data
	values := viewmodel{
		"enrolaction":  "/enrol/submit",
		"disabled":     template.HTMLAttr("disabled='true'"),
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
	view.Views().ExecuteTemplate(w, "work_enrol", values)
}

func SubmitEnrol(w http.ResponseWriter, r *http.Request) {
	if err := onlyPostAllowed(w, r); err != nil {
		return
	}
	app, err := fetchApplicant(w, r, "applicantid")
	if err != nil {
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
	app.Data.OrtSum = ortint(html.EscapeString(r.PostFormValue("ort")))
	app.Data.OrtMath = ortint(html.EscapeString(r.PostFormValue("ortmath")))
	app.Data.OrtPhys = ortint(html.EscapeString(r.PostFormValue("ortphys")))
	app.Data.SchoolOk = r.PostFormValue("schoolok") != ""
	app.Data.OblastOk = r.PostFormValue("oblastok") != ""
	app.Data.OrtOk = r.PostFormValue("ortok") != ""

	oblastID := oblastid(html.EscapeString(r.PostFormValue("district")))
	if app.Data.Enrolled.IsZero() {
		app.Data.Enrolled = time.Now()
	}
	if app.Data.OblastID != oblastID {
		var o model.Oblast
		model.Db().First(&o, oblastID)
		app.Data.Oblast = o
	}
	model.Db().Save(&app)
	view.Views().ExecuteTemplate(w, "work", "")
}

func onlyPostAllowed(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return errors.New("Message not allowed: " + r.Method)
	}
	return nil
}

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
