package controller

import (
	"net/http"
	"github.com/geobe/gostip/go/model"

	"encoding/json"
	"strconv"

	"fmt"
)

func FindByYear(w http.ResponseWriter, r *http.Request){
	year := atoint(r.FormValue("year"))

	var data model.ExamReference

	db := model.Db()

	db.First(&data, "year = ?", year)

	toSend := make([]int, model.NQESTION+1)

	toSend[0] = model.NQESTION
	for i,_ := range toSend{
		if i>0 {
			toSend[i] = data.Results[i - 1]
		}
	}

	a,_ := json.Marshal(toSend)

	w.Write(a)
}

func SubmitExamRef(w http.ResponseWriter, r *http.Request)  {
	var xref model.ExamReference
	year := atoint(r.FormValue("year"))
	db := model.Db()
	db.First(&xref," year = ?", year)
	xref.Year = year
	for i,_ :=range xref.Results {
		task := "task"+strconv.Itoa(i+1)
		xref.Results[i] = atoint(r.FormValue(task))
	}


	fmt.Println(year,xref.Results)
	db.Save(&xref)

	http.Redirect(w,r,"work",301)
}