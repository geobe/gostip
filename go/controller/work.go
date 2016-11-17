package controller

import (
	"fmt"
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html"
	"net/http"
	"strconv"
)

func HandleWork(w http.ResponseWriter, r *http.Request) {
	session, _ := SessionStore().Get(r, S_DKFAI)
	view.Views().ExecuteTemplate(w, "work",
		fmt.Sprintf("user: %s, ID: %s, Name: %s",
			session.Values["fullname"],
			session.ID, session.Name()))
}

func Find(w http.ResponseWriter, r *http.Request) {
	//session, _ := SessionStore().Get(r, S_DKFAI)
	r.ParseForm()
	lastName := html.EscapeString(r.PostFormValue("lastname"))
	firstName := html.EscapeString(r.PostFormValue("firstname"))
	applicants := findApplicants(lastName, firstName)
	//fmt.Printf("found %d applicants\n", len(applicants))
	//for i, ap := range applicants {
	//	fmt.Printf("Applicant %d: %s %s\n", i, ap.Data.FirstName, ap.Data.LastName)
	//}
	view.Views().ExecuteTemplate(w, "qresult", applicants)
}

func Enrol(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	appId, err := strconv.Atoi(html.EscapeString(r.PostFormValue("appid")))
	if err != nil {
		http.Error(w, "Get error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var app model.Applicant
	db := model.Db()
	db.Preload("Data").Preload("Data.Oblast").First(&app, appId)
}

func findApplicants(ln, fn string) (apps []model.Applicant) {
	db := model.Db()
	db.Preload("Data").
		Joins("INNER JOIN applicant_data ON applicants.id = applicant_data.applicant_id").
		Where("applicant_data.deleted_at IS NULL").
		Where("applicant_data.last_name like ?", ln+"%").
		Where("applicant_data.first_name like ?", fn+"%").
		Find(&apps)
	return
}
