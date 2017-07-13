package model

import (
	"strconv"
	"strings"
	"time"

)

// the exam tasks and maximal results in a specific year's test.
// Results array records the maximal points achievable in a question.
// Zero values denote nonexisting questions
type ExamReference struct {
	Model
	Year       int
	Results    [NQESTION]int `gorm:"-"`
	Resultsave string
}

// to easily store grant exam reference values in db with gorm,
// an array of (transient) integer result values is converted to a
// string before saving to db.
// To maintain a short history of changes of exam reference values,
// tracebility is supported by recording timestamp and login of each
// user who initiated the update and last updaters id.
func (xref *ExamReference) BeforeSave() (err error) {
	var r string
	for _, val := range xref.Results {
		r = r + strconv.Itoa(val) + " "
	}
	xref.Resultsave = r
	upid, upsig := signature()
	xref.UpdatedBy += upsig + "(" + time.Now().Format("02.01.2006 15:04:05") + ") "
	xref.Updater = upid
	err = nil
	return
}

// after loading struct fields from db, exam reference values
// are converted back into an array of int
func (xref *ExamReference) AfterFind() (err error) {
	var conv [NQESTION]int
	saved := strings.Split(strings.TrimSpace(xref.Resultsave), " ")
	err = nil
	for i, val := range saved {
		if i >= NQESTION {
			break
		}
		conv[i], err = strconv.Atoi(val)
		if err != nil {
			break
		}
	}
	xref.Results = conv
	return
}
