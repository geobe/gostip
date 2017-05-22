package main

import (
	"github.com/geobe/gostip/go/controller"
	"github.com/geobe/gostip/go/model"
	"net/http"
	"log"
)

func main() {
	// prepare database
	model.Setup("")
	db := model.Db()
	defer db.Close()
	model.ClearTestDb(db)
	model.InitTestDb(db)

	mux := controller.SetRouting()

	// konfiguriere server
	server := &http.Server{
		//Addr: "127.0.0.1:8090",
		Addr:    "0.0.0.0:8090",
		Handler: mux,
	}
	// und starte ihn
	server.ListenAndServe()
	log.Printf("server started\n")
}
