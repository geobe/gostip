package main

import (
	"testing"
	"github.com/geobe/gostip/go/model"
	"fmt"
	"reflect"
	"github.com/pkg/errors"
	"github.com/geobe/gostip/go/controller"
)

func TestStaleUpdate(t *testing.T) {
	// prepare database
	model.Setup("")
	db := model.Db()
	defer db.Close()
	model.ClearTestDb(db)
	model.InitTestDb(db)

	var app0, app1, app2 model.Applicant

	db.Preload("Data").Preload("Data.Oblast").First(&app0)
	db.Preload("Data").Preload("Data.Oblast").First(&app1)
	db.Preload("Data").Preload("Data.Oblast").First(&app2)

	fn := app1.Data.FirstName

	// two way update
	app1.Data.FirstName = fn + "_der_Erste"
	app2.Data.FirstName = fn + "_der_Zweite"
	// mine is updated, but not other, -> no conflict
	app1.Data.Email = "blah@blubb.com"
	// other updated, I did not change -> could be automerged
	app2.Data.LastName = "Hotzenplotz"

	db.LogMode(false).Save(&app1)
	for _, e := range db.LogMode(false).Save(&app2).GetErrors() {
		fmt.Printf("save app2 error: \"%v\"\n", e)
	}
	d, err := getDiff(app1.Data, app2.Data)
	if err != nil {
		fmt.Printf("error in differencing %v\n", err)
	} else {
		for k, v := range d {
			fmt.Printf("diff [%s]: %v <-> %v\n", k, v[0], v[1])
		}
	}

	fmt.Println()

	merge, err := controller.MergeDiff(app0.Data, app1.Data, app2.Data, true)
	if err != nil {
		fmt.Printf("error in merging %v\n", err)
	} else {
		for k, v := range merge {
			fmt.Printf("merge auto [%s]: %v <-> %v is conflict: %v\n", k, v.Mine, v.Other, v.Conflict)
		}
	}

	fmt.Println()

	merge, err = controller.MergeDiff(app0.Data, app1.Data, app2.Data, false)
	if err != nil {
		fmt.Printf("error in merging %v\n", err)
	} else {
		for k, v := range merge {
			fmt.Printf("merge show [%s]: %v <-> %v is conflict: %v\n", k, v.Mine, v.Other, v.Conflict)
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
