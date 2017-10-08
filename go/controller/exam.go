package controller

import (
	"net/http"
	"github.com/geobe/gostip/go/model"
	"strconv"
	"fmt"
	"github.com/geobe/gostip/go/view"
)

func ExamRef(w http.ResponseWriter, r *http.Request){
	var examrefs []*model.ExamReference = make([]*model.ExamReference,0)
	resultsCounter := make([]int,model.NQESTION)

	for i,_ := range resultsCounter{
		resultsCounter[i] = i+1;
	}
	db := model.Db()

	db.Find(&examrefs)
	values := viewmodel{
		"Examref" : examrefs,
		"Rescounter" : resultsCounter,
		"NQ" : model.NQESTION,
	}

	view.Views().ExecuteTemplate(w, "adminarea",values)
}

func SubmitExamRef(w http.ResponseWriter, r *http.Request)  {
	var xref model.ExamReference
	year := atoint(r.FormValue("year"))
	db := model.Db()
	db.First(&xref," year = ?", year)
	xref.Year = year
	for i,_ :=range xref.Results {
		task := "mr"+strconv.Itoa(i+1)
		xref.Results[i] = atoint(r.FormValue(task))
	}


	fmt.Println(year,xref.Results)
	db.Save(&xref)

	http.Redirect(w,r,"work",301)
}