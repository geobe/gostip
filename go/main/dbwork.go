package main

import (
	"fmt"
	"github.com/geobe/gostip/go/model"
	"time"
)

func main1() {
	model.Setup("")
	db := model.Db()
	//db.LogMode(true)
	defer db.Close()

	//model.ClearTestDb(db)
	model.InitTestDb(db)

	// retrieve an applicant from db
	var appli model.Applicant
	db.Preload("Data").Preload("Data.Oblast").First(&appli)

	// check load
	fmt.Printf("%s %s aus %s (%d-> %s) erreicht %v\n",
		appli.Data.FirstName, appli.Data.LastName,
		appli.Data.Home, appli.Data.OblastID, appli.Data.Oblast.Name, appli.Data.Results)

	// change some data and save again
	appli.Data.Results = [model.NQESTION]int{7, 7, 7, 7, 7, 7, 7, 7, 7, 8}
	db.Save(&appli)
	// wait a bit
	time.Sleep(time.Second)
	// change back and save again
	appli.Data.Results = [model.NQESTION]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	db.Save(&appli)

	// find applicants from applicant data with different where conditions
	var applicants []model.Applicant
	qs := "Du%"

	db.Preload("Data").
		Joins("JOIN applicant_data ON applicants.id = applicant_data.applicant_id "+
			"AND applicant_data.last_name like ?", qs).
		Where("applicant_data.deleted_at IS NULL").
		Find(&applicants)
	for i, ap := range applicants {
		fmt.Printf("Applicant %d: %s %s\n", i, ap.Data.FirstName, ap.Data.LastName)
	}

	fmt.Println()

	qs = "%ns%"
	db.Preload("Data").
		Joins("INNER JOIN applicant_data ON applicants.id = applicant_data.applicant_id").
		Where("applicant_data.deleted_at IS NULL").
		Where("applicant_data.last_name like ?", qs).
		Find(&applicants)
	for i, ap := range applicants {
		fmt.Printf("Applicant %d: %s %s\n", i, ap.Data.FirstName, ap.Data.LastName)
	}

	fmt.Println("")

	qs = "%"
	db.Preload("Data").Preload("Data.Oblast").
		Joins("JOIN applicant_data ON applicants.id = applicant_data.applicant_id "+
			"JOIN oblasts ON oblasts.id = applicant_data.oblast_id").
		Where("applicant_data.deleted_at IS NULL").
		Where("oblasts.name like ?", qs).Order("applicant_data.ort_sum desc").
		Find(&applicants)
	for i, ap := range applicants {
		fmt.Printf("Applicant %d: %s %s from %s\n", i, ap.Data.FirstName, ap.Data.LastName, ap.Data.Oblast.Name)
	}

	fmt.Println("")
}
