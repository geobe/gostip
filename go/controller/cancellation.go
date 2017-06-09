package controller

import (
	"github.com/geobe/gostip/go/view"
	"html"
	"net/http"
	"github.com/geobe/gostip/go/model"
	"time"
	"github.com/justinas/nosurf"
)

// ShowCancellation is handler to show the selected applicant from the
// search select element for cancellation tab. It returns
// an html page fragment that is inserted into the respective tab area.
func ShowCancellation(w http.ResponseWriter, r *http.Request) {
	if checkMethodAllowed(http.MethodPost, w, r) != nil {
		return
	}
	if err := parseSubmission(w, r); err != nil {
		return
	}
	flag := html.EscapeString(r.PostFormValue("flag"))
	app, err := fetchApplicant(w, r, "appid", flag != "")
	if err != nil {
		return
	}

	values := viewmodel{
		"csrftoken": nosurf.Token(r),
		"csrfid": "csrf_id_cancel",
		"language":     view.PreferedLanguages(r) [0],
	}
	setViewModel(app, values)
	view.Views().ExecuteTemplate(w, "work_cancellation", values)
}

// SubmitCancellation is handler that accepts form submissions from the cancellation tab.
// Only http POST method is accepted.
func SubmitCancelation(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		processCancellation(w, r, true)
	}
}

// processCancellation deletes or undeletes an applicant in the database
func processCancellation(w http.ResponseWriter, r *http.Request, enrol bool) {
	if err := parseSubmission(w, r); err != nil {
		return
	}
	flag := html.EscapeString(r.PostFormValue("flag"))
	undo := flag != ""
	action := html.EscapeString(r.PostFormValue("action"))
	app, err := applicantFromSession(action, r) // fetchApplicant(w, r, "appid", action, undo)
	if err == nil {
		if undo {
			app.DeletedAt = nil
			app.Data.CancelledAt = time.Time{}
			model.Db().Unscoped().Save(&app)
		} else {
			tx := model.Db().Begin()
			app.Data.CancelledAt = time.Now()
			tx.Save(&app)
			tx.Delete(&app)
			tx.Commit()
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
