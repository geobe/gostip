package main

import (
	"testing"
	"github.com/geobe/gostip/go/model"
	"reflect"
	"fmt"
)

func TestReflect(t *testing.T) {
	var a model.ApplicantData = model.ApplicantData{}
	var pa *model.ApplicantData = &model.ApplicantData{}

	refA := reflect.ValueOf(a)
	refPa := reflect.ValueOf(pa).Elem()

	fmt.Printf("Type %s\n", refA.Type().Name())
	for i := 0; i < refA.NumField(); i++ {
		fieldInfo := refA.Type().Field(i)
		fieldName := fieldInfo.Name//refA.Field(i)
		fmt.Printf("\tField %s has type %s\n", fieldName, fieldInfo.Type.Name())
	}

	fmt.Printf("Type %s\n", refPa.Type().Name())
	for i := 0; i < refPa.NumField(); i++ {
		fieldInfo := refPa.Type().Field(i)
		fieldName := fieldInfo.Name//refPa.Field(i)
		fmt.Printf("\tField %s has type %s\n", fieldName, fieldInfo.Type.Name())
	}
}
