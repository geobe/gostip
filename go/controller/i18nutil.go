package controller

import (
	"github.com/geobe/gostip/go/model"
	"os"
	"log"
	"encoding/csv"
)

func ReadCsv(file string) map[string]map[string]string {
	var translations = make(map[string]map[string]string)

	configfile := model.Base + "/config/" + file
	if rd, err := os.OpenFile(configfile, os.O_RDONLY, 0666); err == nil {
		csvReader := csv.NewReader(rd)
		if all, er2 := csvReader.ReadAll(); er2 == nil {
			head := all[0]
			for _, lang := range head[1:] {
				translations[lang] = make(map[string]string)
			}
			body := all[1:]
			for _, line := range body {
				key := line[0]
				for i, val := range line[1:] {
					translations[head[i+1]][key] = val
				}
			}
		} else {
			log.Printf ("read error %x", er2)
		}

	}else {
		log.Printf("open error %x", err)
	}
	return translations
}
