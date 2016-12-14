package model

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
	"time"
	"fmt"
)

// # of questions that could appear in one test
const NQESTION = 10

// a type to classify language
type Lang int

// the languages relevant for the test
const (
	NONE Lang = iota
	DE
	EN
	RU
	KG
	FR
	ES
	CN
	OTHER = 99
)

// an applicant for a place at university
type Applicant struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time    `sql:"index"`
	Data      ApplicantData // `gorm:"ForeignKey:DataID"`
}

// all data of an applicant are stored in a separate structure in order
// to maintain a full history of changes to these sensitrive data
// form: tags are used to identify html form fields and request parameters
type ApplicantData struct {
	Model
	ApplicantID    uint
	Number         uint `gorm:"AUTO_INCREMENT" form:"applid"`
	LastName       string `form:"lastname"`
	FirstName      string `form:"firstname"`
	FathersName    string `form:"fathersname"`
	Phone          string `form:"phone"`
	Email          string `form:"email"`
	Home           string `form:"home"`
	School         string `form:"school"`
	SchoolOk       bool `form:"schoolok"`
	Oblast         Oblast                   // Belongs To Association
	OblastID       uint `form:"district"`
	OblastOk       bool `form:"districtok"`
	OrtSum         int16 `form:"ort"`
	OrtMath        int16 `form:"ortmath"`
	OrtPhys        int16 `form:"ortphys"`
	OrtOk          bool `form:"ortok"`
	Results        [NQESTION]int `gorm:"-" form:"r#1"` // marks multiplied by 10
	Resultsave     string
	LanguageResult int `form:"languageresult"`
	Language       Lang `form:"language"`
	EnrolledAt     time.Time `form:"enrolledat"`
	CancelledAt    time.Time `form:"cancelledat"`
}

// a district in Kyrgyzstan
type Oblast struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

var InitialOblasts = []Oblast{
	{Name: "Bishkek City"},
	{Name: "Osh City"},
	{Name: "Batken"},
	{Name: "Chuy"},
	{Name: "Jalal-Abad"},
	{Name: "Naryn"},
	{Name: "Osh"},
	{Name: "Talas"},
	{Name: "Yssykköl"},
	{Name: "Foreign"},
}

var InitialLanguages = map[Lang]string{
	NONE:        "keine",
	DE:        "Deutsch",
	EN:        "Englisch",
	RU:        "Russisch",
	KG:        "Kirgisisch",
	FR:        "Französisch",
	ES:        "Spanisch",
	CN:        "Chinesisch",
	OTHER:        "andere",
}

// to easily store grant exam results in db with gorm,
// an array of (transient) integer result values is converted to a
// string before saving to db
func (appdata *ApplicantData) BeforeSave() (err error) {
	var r string
	for _, val := range appdata.Results {
		r = r + strconv.Itoa(val) + " "
	}
	appdata.Resultsave = r
	err = nil
	return
}

// after loading struct fields from db, exam result values
// are converted back into an array of int
func (appdata *ApplicantData) AfterFind() (err error) {
	var conv [NQESTION]int
	saved := strings.Split(strings.TrimSpace(appdata.Resultsave), " ")
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
	appdata.Results = conv
	return
}

// identify current user
func signature() (uint, string) {
	return 42, "Me@" + time.Now().Format("01.02.2006-03:04:05.9999")
}

// to maintain a full history of changes of applicant data,
// old data are kept in the db using gorms 'DeletedAt' mechanism
// and a new record with updated data is saved to db. Tracebility
// is ensured by recording the identity of the user who
// initiated the update.
// concurrent conflicting updates are detected by comparing
// changed_at field
func (app *Applicant) BeforeUpdate(tx *gorm.DB) (err error) {
	var appdb Applicant
	tx.First(&appdb, app.ID)
	fmt.Printf("updating applicant upd@%v onto upd@%v\n", app.UpdatedAt, appdb.UpdatedAt)
	if appdb.UpdatedAt != app.UpdatedAt {
		fmt.Printf("stale error updating applicant %d\n", app.ID)
		err = fmt.Errorf("stale object")
		return
	}
	data := app.Data
	upid, upsig := signature()
	data.Model = Model{UpdatedBy: upsig, Updater: upid}
	tx.Delete(&app.Data)
	app.Data = data
	err = nil
	return
}
