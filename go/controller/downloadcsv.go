package controller

import (
	"os"
	"encoding/csv"
	"github.com/geobe/gostip/go/model"
	"strconv"
	"log"
	"net/http"
)

var initdata = make([]string,100)
var first = []string{"Test_ID", "Vorname", "Nachname", "VornameTx", "NachnameTx", "Telefon", "Email", "Oblast", "Stadt", "Ort", "ortMath", "ortPhys"}

func DownloadCsv(w http.ResponseWriter,r *http.Request) {

	copy(initdata, first)
	var database []*model.ApplicantData = make([]*model.ApplicantData,0)

	db := model.Db()
	db.Find(&database)

	i :=0
	for _,value := range database{
		for _,ch := range value.Resultsave{
			if ch == 32 {
				initdata[12+i] = strconv.Itoa(i+1)
				i++
			}
		}
		break
	}

	initdata[12+i] = "Gesamt"
	initdata[13+i] = "D/E"
	initdata[14+i] = "Language"
	var data  = make([]string, 15+i)
	copy(data,initdata)
	file, err := os.Create("src/github.com/geobe/gostip/resources/list/resultlist.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	error := writer.Write(data)
	checkError("Cannot write to file", error)


	for _,value := range database{
		data[0] = strconv.Itoa(int(value.ApplicantID))
		data[1] = value.FirstName
		data[2] = value.LastName
		data[3] = value.FirstNameTx
		data[4] = value.LastNameTx
		data[5] = value.Phone
		data[6] = value.Email
		data[7] = value.Oblast.Name
		data[8] = value.Home
		data[9] = strconv.Itoa(int(value.OrtSum))
		data[10] = strconv.Itoa(int(value.OrtMath))
		data[11] = strconv.Itoa(int(value.OrtPhys))

		results := value.Resultsave
		sum := 0
		numb :=0
		i = 0
		for _,sub := range results{
			if sub == 32 {
				data[12+i] = strconv.Itoa(numb)
				i++
				sum += numb
				numb = 0
				continue
			}
			numb += int(sub-48)*10
		}
		data[12+i] = strconv.Itoa(sum)
		data[13+i] = model.InitialLanguages[value.Language]
		data[14+i] = strconv.Itoa(value.LanguageResult)
		err := writer.Write(data)
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}


