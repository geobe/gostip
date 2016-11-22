// work.go implements all controller actions and their helper functions
// for the superordinated work.html template. This also comprises the
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
)

// load the "work" template. Tabs are included depending on role of current user.
func HandleWork(w http.ResponseWriter, r *http.Request) {
	roles := viewmodel{}
	addRoles(r, roles)
	view.Views().ExecuteTemplate(w, "work", roles)
}

// handler function that executes a database search for applicants and
// returns an html fragment for a select box.
func Find(w http.ResponseWriter, r *http.Request) {
	//session, _ := SessionStore().Get(r, S_DKFAI)
	r.ParseForm()
	lastName := html.EscapeString(r.PostFormValue("lastname"))
	firstName := html.EscapeString(r.PostFormValue("firstname"))
	action := html.EscapeString(r.PostFormValue("action"))
	enrol := action == "enrol"
	applicants := findApplicants(lastName, firstName, enrol)
	fmt.Printf("found %d applicants\n", len(applicants))
	for i, ap := range applicants {
		fmt.Printf("Applicant %d: %s %s %v\n", i, ap.Data.FirstName, ap.Data.LastName, ap.Data.EnrolledAt)
	}
	view.Views().ExecuteTemplate(w, "qresult", applicants)
}

// find applicants based on lastname and/or firstname. Per default, a wildcard search
// character (%) is appended to the search strings and query uses LIKE condition.
func findApplicants(ln, fn string, enrol bool) (apps []model.Applicant) {
	var qs string
	if enrol {
		qs = "applicant_data.enrolled_at <'1900-01-01'"
	} else {
		qs = "applicant_data.enrolled_at > '" +
			time.Now().Format("2006") + "-01-01'"
	}
	db := model.Db()
	db.Preload("Data").
		Joins("INNER JOIN applicant_data ON applicants.id = applicant_data.applicant_id").
		Where("applicant_data.deleted_at IS NULL").
		Where(qs).
		Where("applicant_data.last_name like ?", ln+"%").
		Where("applicant_data.first_name like ?", fn+"%").
		Find(&apps)
	return
}

// add user role fields to the viewmodel map according to the role privileges of current user
func addRoles(r *http.Request, data viewmodel) (err error) {
	session, err := SessionStore().Get(r, S_DKFAI)
	if err != nil {
		return
	}
	role, ok := session.Values["role"].(int)
	fmt.Printf("role = %d\n", role)
	if !ok {
		err = errors.New("no role defined")
		return
	}
	if role&model.U_ANY != 0 {
		data["authany"] = true
	}
	if role&model.U_ENROL != 0 {
		data["authenrol"] = true
	}
	if role&model.U_POFF != 0 {
		data["authpoff"] = true
	}
	if role&model.U_UADMIN != 0 {
		data["authuadmin"] = true
	}
	if role&model.U_FULLADMIN != 0 {
		data["authfulladmin"] = true
	}
	if role&model.U_ALL != 0 {
		data["authall"] = true
	}
	return
}
