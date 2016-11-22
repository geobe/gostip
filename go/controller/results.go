// results.go implements handler and helper functions
// for the results and resultslist tabs
package controller

import "github.com/geobe/gostip/go/model"

// add a slice of possible test results for the given year to the viewmodel
func addResultsConfig(y int, data viewmodel) {
	var exref model.ExamReference
	model.Db().First(&exref, "year = ?", y)
	if exref.ID == 0 {
		return
	}
	var rslice []int
	for i, v := range exref.Results {
		if v == 0 || i == model.NQESTION-1 {
			rslice = exref.Results[0:i]
			break
		}
	}
	data["results"] = rslice
}
