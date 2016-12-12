// a couple of predefined data for development time
package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

var classes = []interface{}{
	Applicant{},
	ApplicantData{},
	Oblast{},
	User{},
	ExamReference{},
}

// create some test applicants in a developpment db
func InitTestDb(db *gorm.DB) *gorm.DB {

	// Migrate the schema
	db.AutoMigrate(classes...)

	// initialize if db empty
	var nOblasts, nappl, nxref int
	if db.Model(&Oblast{}).Count(&nOblasts); nOblasts < 1 {
		fmt.Print("No oblasts in db, creating new\n")
		for _, oblast := range InitialOblasts {
			db.Create(&oblast)
		}
	}
	// retrieve oblasts from db
	var oblasts []Oblast
	db.Find(&oblasts)

	if db.Model(&ExamReference{}).Count(&nxref); nxref < 1 {
		xref := ExamReference{
			Year:    time.Now().Year(),
			Results: [NQESTION]int{50, 50, 20, 30, 100, 150, 100, 100, 0, 0},
		}
		db.Create(&xref)
	}

	if db.Find(&Applicant{}).Count(&nappl); nappl < 1 {
		data := ApplicantData{
			LastName:    "Gans",
			FirstName:   "Gisbert",
			FathersName: "Gisbertovich",
			Phone:       "040441777",
			Home:        "Ducksburg",
			School:      "Nr. 14",
			Email:       "Giga@goosemail.com",
			Oblast:      oblasts[7],
			OrtSum:      111,
			OrtMath:     66,
			OrtPhys:     33,
			Results:     [NQESTION]int{10, 20, 30, 40, 555, 60, 75, 80, 9, 10},
		}
		goose := Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Gans",
			FirstName:   "Franz",
			FathersName: "Fietjevich",
			Phone:       "046774417",
			Home:        "Franzhausen",
			School:      "Nr. 66",
			Email:       "FraGa@goosemail.com",
			Oblast:      oblasts[6],
			OrtSum:      100,
			OrtMath:     66,
			OrtPhys:     33,
			Results:     [NQESTION]int{10, 25, 35, 40, 50, 25, 70, 80, 90, 100},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Gans",
			FirstName:   "Gertrude",
			FathersName: "Gaggova",
			Phone:       "0467734567",
			Home:        "Ducksburg",
			School:      "Nr. 35",
			Email:       "Gerti@goosemail.com",
			Oblast:      oblasts[6],
			OrtSum:      178,
			OrtMath:     99,
			OrtPhys:     36,
			Results:     [NQESTION]int{10, 25, 35, 40, 50, 25, 70, 80, 90, 100},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Duck",
			FirstName:   "Daisy",
			FathersName: "Waltova",
			Phone:       "08980080099",
			Home:        "München",
			School:      "Nr. 5",
			Email:       "Daisy@disney.com",
			Oblast:      oblasts[2],
			OrtSum:      230,
			OrtMath:     100,
			OrtPhys:     66,
			Results:     [NQESTION]int{10, 25, 35, 40, 50, 25, 70, 80, 90, 100},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Duck",
			FirstName:   "Donald",
			FathersName: "Disneyvich",
			Phone:       "08955443322",
			Home:        "Oberpfaffenhofen",
			School:      "Nr. 19",
			Email:       "Donald@disney.com",
			Oblast:      oblasts[2],
			OrtSum:      101,
			OrtMath:     50,
			Results:     [NQESTION]int{10, 25, 35, 40, 50, 25, 70, 80, 90, 100},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Quack",
			FirstName:   "Primus",
			FathersName: "Disneyvich",
			Phone:       "04055443322",
			Home:        "Hamburg",
			School:      "Nr. 1",
			Email:       "Primus.Quack@disney.com",
			Oblast:      oblasts[5],
			OrtSum:      245,
			OrtMath:     100,
			OrtPhys:     100,
			EnrolledAt:  time.Now(),
			Results:     [NQESTION]int{20, 30, 0, 55, 45, 75, 100, 100},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
	}
	var nuser int
	if db.Model(&User{}).Count(&nuser); nuser < 1 {
		users := InitialUsers()
		for _, user := range users {
			db.Save(&user)
		}
	}
	return db
}

// clear development db
func ClearTestDb(db *gorm.DB) {

	for _, class := range classes {
		db.Unscoped().Delete(class)
		//db.DropTableIfExists(class)
	}

}
