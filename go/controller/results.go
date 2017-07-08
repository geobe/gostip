package controller

import (
	"github.com/geobe/gostip/go/model"
	"github.com/geobe/gostip/go/view"
	"html/template"
	"net/http"
	"time"
	"fmt"
	"github.com/justinas/nosurf"
	"log"
	"strings"
	"strconv"
)

// controller function to show applicant data for results editing
func ShowResults(w http.ResponseWriter, r *http.Request) {
	if checkMethodAllowed(http.MethodPost, w, r) != nil {
		return
	}
	app, err := fetchApplicant(w, r, "appid")
	//fmt.Printf("got applicant %s\n", app.Data.FirstName)
	if err != nil {
		return
	}
	i18nlanguage := view.PreferedLanguages(r) [0]
	values := viewmodel{
		"disabled": template.HTMLAttr("disabled='true'"),
		"oblasts":  view.OblastsI18n(i18nlanguage),
		"csrftoken": nosurf.Token(r),
		"csrfid": "csrf_id_results",
		"language":     i18nlanguage,
		"languages":        view.LanguagesI18n(i18nlanguage),
	}
	setViewModel(app, values)
	addResultsConfig(time.Now().Year(), app, values)
	view.Views().ExecuteTemplate(w, "work_results", values)

}

// submit results into ApplicantData struct
func SubmitResults(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodPost, w, r); err == nil {
		saveApplicantSubmission(w, r, true)
	}
}

// controller function to prepare a csv file with applicant data for download
func GetResultsCsv(w http.ResponseWriter, r *http.Request) {
	if err := checkMethodAllowed(http.MethodGet, w, r); err == nil {
		getApplicantCsv(w, r)
	}
}

func getApplicantCsv(w http.ResponseWriter, r *http.Request) {
	var year int
	ok := setIfPosted(&year, "year", r)
	if ! ok {
		year = time.Now().Year()
	}
	var exref model.ExamReference
	model.Db().First(&exref, "year = ?", year)
	if exref.ID == 0 {
		http.Error(w, "No configuration data for year ", http.StatusNotFound)
		log.Printf("No configuration data for year %v, status %v\n", year, http.StatusNotFound)
		return
	}
	nquestions := exref.QuestionsCount()
	noresult := strings.Repeat("0 ", model.NQESTION)
	var participants []model.ApplicantData
	i18nlanguage := view.PreferedLanguages(r) [0]

	t0 := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	t1 := time.Date(year + 1, time.January, 1, 0, 0, 0, 0, time.UTC)
	model.Db().Order("last_name").Preload("Oblast").
		Find(&participants, "resultsave <> ? AND enrolled_at >= ? AND enrolled_at < ? ", noresult, t0, t1)
	fmt.Fprint(w, makeHeaderLines(nquestions))
	for _, app := range participants {
		fmt.Fprintf(w, "%d;%s;%s;%s;%s;\"%s\";%s;%s;%s;%s;%d;%d;%d",
			app.ID,
			app.LastName, app.LastNameTx, app.FirstName, app.FirstNameTx,
			app.Phone, app.Email, app.Home, app.School,
			view.I18n(app.Oblast.Name, i18nlanguage),
			app.OrtSum, app.OrtMath, app.OrtPhys)
		for i := 0; i < nquestions; i++ {
			res := float32(app.Results[i]) / 10.
			fmt.Fprintf(w, ";%.1f", res)
		}
		fmt.Fprintf(w, ";;%s;%.1f\n", app.Language.Lang(), float32(app.LanguageResult) / 10.)
	}
}
func makeHeaderLines(nquestion int) string {
	h11 := "ID;Name;NameTX;Vorname;VornameTx;" +
		"Telefon;Email;Stadt;Schule;Oblast;" +
		"ORT gesamt;Mathe;Physik;" +
		"Mathematisch-logische Aufgaben;"
	h12 := "Sprache;\n"
	h21 := ";;;;;" +
		";;;;;" +
		";;;"
	h22 := "Gesamt;d/e/-;;ORT (35%);Mat.-log. (45%);Sprache (20%);Gesamt;Bemerkung;Gebühr;Lebenshaltung\n"
	qn := ""
	pl := ""
	for i := 1; i <= nquestion; i++ {
		qn = qn + strconv.Itoa(i) + ";"
		pl = pl + ";"
	}
	return h11 + pl + h12 + h21 + qn + h22
}

// add a slice of possible test results for the given year to the viewmodel
func addResultsConfig(y int, app model.Applicant, data viewmodel) {
	var exref model.ExamReference
	model.Db().First(&exref, "year = ?", y)
	if exref.ID == 0 {
		return
	}
	var nq int
	var results [model.NQESTION]map[string]float32

	for i, v := range exref.Results {
		if v == 0 || i == model.NQESTION - 1 {
			nq = i
			break
		}
	}

	for i := 0; i <= nq; i++ {
		results[i] = map[string]float32{
			"val": float32(app.Data.Results[i]) / 10.,
			"max": float32(exref.Results[i]) / 10.,
			"no":  float32(i + 1),
		}
	}
	data["results"] = results[:nq]
	data["languageresult"] = fmt.Sprintf("%.1f", float32(app.Data.LanguageResult) / 10.)
	data["lang"] = app.Data.Language
}
