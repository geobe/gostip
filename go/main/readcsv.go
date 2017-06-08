package main

import (
	"github.com/geobe/gostip/go/controller"
	"log"
)

func main() {
	trmap:=controller.ReadCsv("allInOne.csv")
	log.Print(trmap["de"]["_hello_personal"])
	log.Print(trmap["ru"]["_hello_personal"])
}
