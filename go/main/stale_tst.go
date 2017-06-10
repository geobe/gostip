package main

import (
	"github.com/geobe/gostip/go/model"
	"fmt"
	"reflect"
	"github.com/pkg/errors"
	"github.com/geobe/gostip/go/controller"
)

//func TestStaleUpdate(t *testing.T) {
func main() {
	// prepare database
	model.Setup("")
	db := model.Db()
	defer db.Close()
	model.ClearTestDb(db)
	//model.InitTestDb(db)

	var app0, app1, app2 model.Applicant

	db.Preload("Data").Preload("Data.Oblast").First(&app0)
	db.Preload("Data").Preload("Data.Oblast").First(&app1)
	db.Preload("Data").Preload("Data.Oblast").First(&app2)

	var naryn, batken model.Oblast

	db.Where("name = ?", "Naryn").First(&naryn)
	db.Where("name = ?", "Batken").First(&batken)

	fmt.Printf("Oblast %v\n", naryn)

	fn := app1.Data.FirstName
	ln := app1.Data.LastName
	pn := app1.Data.FathersName

	// two way update
	app1.Data.FirstName = fn + "_1"
	app1.Data.LastName = ln + "bert"
	app1.Data.Oblast = batken
	app1.Data.OblastID = batken.ID
	app1.Data.Results[2]++
	app1.Data.Results[3]++
	app1.Data.Results[4]++
	app2.Data.Oblast = naryn
	app2.Data.OblastID = naryn.ID
	app2.Data.Results[3]++
	app2.Data.Results[4]--
	app2.Data.Results[5]++
	app2.Data.FirstName = fn + "_2"
	app2.Data.LastName = ln + "bert"
	app2.Data.FathersName = pn + "hurz"
	// mine is updated, but not other, -> no conflict
	app1.Data.Email = "blah@blubb.com"
	// other updated, I did not change -> could be automerged
	//app2.Data.LastName = "Hotzenplotz"

	db.LogMode(false).Save(&app1)
	for _, e := range db.LogMode(false).Save(&app2).GetErrors() {
		fmt.Printf("save app2 error: \"%v\"\n", e)
	}

	fmt.Printf("before merge\n")
	fmt.Printf("app0.Data: %s %s [%s] from %s(%d) results %v\n",
		app0.Data.FirstName, app0.Data.LastName, app0.Data.Email, app0.Data.Oblast.Name,
		app0.Data.Oblast.ID,  app0.Data.Results)
	fmt.Printf("app1.Data: %s %s [%s] from %s(%d) results %v\n",
		app1.Data.FirstName, app1.Data.LastName, app1.Data.Email, app1.Data.Oblast.Name,
		app1.Data.Oblast.ID,  app1.Data.Results)
	fmt.Printf("app2.Data: %s %s [%s] from %s(%d) results %v\n",
		app2.Data.FirstName, app2.Data.LastName, app2.Data.Email, app2.Data.Oblast.Name,
		app2.Data.Oblast.ID, app2.Data.Results)

	merge, err := controller.MergeDiff(&app0.Data, &app1.Data, &app2.Data, true, "form")
	if err != nil {
		fmt.Printf("error in merging %v\n", err)
	} else {
		fmt.Printf("after merge\n")
		fmt.Printf("app1.Data: %s %s [%s] from %s(%d) results %v\n",
			app1.Data.FirstName, app1.Data.LastName, app1.Data.Email, app1.Data.Oblast.Name,
			app1.Data.OblastID, app1.Data.Results)
		for k, v := range merge {
			cnf := "update"
			ic := ""
			switch v.Conflict {
			case controller.MINE:
				cnf = "mine"
			case controller.THEIRS:
				cnf = "theirs"
			case controller.BOTH:
				cnf = "both"
			case controller.SAME:
				cnf = "same"
			}
			if v.Conflict != controller.NONE {
				ic = ">"
			}
			fmt.Printf("merge auto [%s]: %v <-%s %v is %s\n", k, v.Mine, ic, v.Other, cnf)
		}
	}


}

func getDiff(val1, val2 interface{}) (diffs map[string][]interface{}, err error) {
	v1 := reflect.ValueOf(val1)
	v2 := reflect.ValueOf(val2)
	diffs = make(map[string][]interface{})

	if v1.Type() != v2.Type() {
		err = errors.New("different types")
		return
	}
	for i := 0; i < v1.NumField(); i++ {
		fieldInfo := v1.Type().Field(i)
		fieldName := fieldInfo.Name
		fieldVal1 := v1.Field(i).Interface()
		fieldVal2 := v2.Field(i).Interface()
		if !fieldInfo.Anonymous && fieldVal1 != fieldVal2 {
			diffs[fieldName] = []interface{}{fieldVal1, fieldVal2, }
		}
	}
	return
}
