// work.go implements all controller actions and their helper functions
// for the parent work.html template. This also comprises the
// search functionality. All user work is done on tabs that are dynamically
// loaded into the work template. These have their own controller files.
package controller

import (
	"fmt"
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"github.com/pkg/errors"
	"html"
	"net/http"
	"time"
	"github.com/justinas/nosurf"
	"github.com/geobe/gostip/go/transcription"
)

// load the "work" template. Tabs are included depending on role of current user.
func HandleWork(w http.ResponseWriter, r *http.Request) {
	roles := viewmodel{}
	roles["csrftoken"] = nosurf.Token(r)
	roles["csrfid"] = "csrf_id_find"
	addRoles(r, roles)
	view.Views().ExecuteTemplate(w, "work", roles)
}

// handler function that executes a database search for applicants and
// returns an html fragment for a select box.
func FindApplicant(w http.ResponseWriter, r *http.Request) {
	l := r.Header["Accept-Language"]
	getKyr := transcription.UsesKyrillic(l)
	r.ParseForm()
	lastName := html.EscapeString(r.PostFormValue("lastname"))
	firstName := html.EscapeString(r.PostFormValue("firstname"))
	action := html.EscapeString(r.PostFormValue("action"))
	flag := html.EscapeString(r.PostFormValue("flag"))
	enrol := action == "enrol"
	active := flag == ""
	applicants := findApplicants(lastName, firstName, enrol, active)
	view.Views().ExecuteTemplate(w, "qresult", applicantResultList(applicants, getKyr))
}

func applicantResultList(appls []model.Applicant, getKyr bool) (res []map[string]string) {
	res = make([]map[string]string, len(appls))
	for i, app := range appls {
		isKyr := transcription.IsKyrgyz(app.Data.LastName)
		if isKyr == getKyr {
			res[i] = map[string]string{
				"id":        fmt.Sprintf("%d", app.ID),
				"lastname":  app.Data.LastName,
				"firstname": app.Data.FirstName,
			}
		} else {
			res[i] = map[string]string{
				"id":        fmt.Sprintf("%d", app.ID),
				"lastname":  app.Data.LastNameTx,
				"firstname": app.Data.FirstNameTx,
			}
		}
	}
	return
}

// find applicants based on lastname and/or firstname. Per default, a wildcard search
// character (%) is appended to the search strings and query uses LIKE condition.
// Search is performed also in transcription fields
// ln, fn:      lastname, firstname search strings
// enrol:       true -> searching from the enrol use case for new applicants
// active:      true -> active applicant, not cancelled
func findApplicants(ln, fn string, enrol bool, active bool) (apps []model.Applicant) {
	var qs string
	if enrol {
		// query for newly registered applicants
		qs = "applicant_data.enrolled_at <'1900-01-01'"
	} else {
		// query for enrolled applicants
		qs = "applicant_data.enrolled_at > '" +
			time.Now().Format("2006") + "-01-01'"
	}
	db := model.Db()
	if active {
		// query for active applicants
		db.Preload("Data").
			Joins("INNER JOIN applicant_data ON applicants.id = applicant_data.applicant_id").
			Where("applicant_data.deleted_at IS NULL").
			Where(qs).
			Where("applicant_data.last_name like ? OR applicant_data.last_name_tx like ?",
			ln + "%", ln + "%").
			Where("applicant_data.first_name like ? OR applicant_data.first_name_tx like ?",
			fn + "%", fn + "%").
			Find(&apps)
	} else {
		// query for deleted applicants
		db.Unscoped().Preload("Data").
			Where("applicants.deleted_at IS NOT NULL").
			Joins("INNER JOIN applicant_data ON applicants.id = applicant_data.applicant_id").
			Where("applicant_data.deleted_at IS NULL").
			Where(qs).
			Where("applicant_data.last_name like ?", ln + "%").
			Where("applicant_data.first_name like ?", fn + "%").
			Find(&apps)
	}
	return
}

// add user role fields to the viewmodel map according to the role privileges of current user
func addRoles(r *http.Request, data viewmodel) (err error) {
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		return
	}
	role, ok := session.Values["role"].(int)
	if !ok {
		err = errors.New("no role defined")
		return
	}
	if role & model.U_ANY != 0 {
		data["authany"] = true
	}
	if role & model.U_ENROL != 0 {
		data["authenrol"] = true
	}
	if role & model.U_POFF != 0 {
		data["authpoff"] = true
	}
	if role & model.U_UADMIN != 0 {
		data["authuadmin"] = true
	}
	if role & model.U_FULLADMIN != 0 {
		data["authfulladmin"] = true
	}
	if role & model.U_ALL != 0 {
		data["authall"] = true
	}
	return
}
