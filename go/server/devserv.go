package main

import (
	"fmt"
	"github.com/geobe/gostip/go/controller"
	"github.com/geobe/gostip/go/model"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
	"os"
)

const pages = "/github.com/geobe/gostip/pages/"
const resources = "/github.com/geobe/gostip/resources/"

func err(writer http.ResponseWriter, request *http.Request) {
	//error := request.Header.Get("Status")
	fmt.Fprintf(writer, "Error with path %s",
		request.URL.Path[1:])
}

func main() {
	// prepare database
	model.Setup("")
	db := model.Db()
	defer db.Close()
	model.ClearTestDb(db)
	model.InitTestDb(db)

	mux := mux.NewRouter()
	// finde Working directory = GOPATH
	docbase, _ := os.Getwd()
	// FileServer ist ein Handler, der dieses Verzeichnis bedient
	files := http.FileServer(http.Dir(docbase + pages))
	// Zugriff auf das Verzeichnis via Präfic /pages/
	mux.PathPrefix("/pages/").Handler(http.StripPrefix("/pages/", files))
	// Zugriff auf die Resourcen-Verzeichnisse mit regular expression
	mux.PathPrefix("/{dir:(css|fonts|js)}/").Handler(http.FileServer(http.Dir(docbase + resources)))
	// Zugriff auf *.htm[l] Dateien im /pages Verzeichnis
	mux.Handle("/{dir:\\w+\\.html?}", http.FileServer(http.Dir(docbase+pages)))
	// error
	mux.HandleFunc("/err", err)
	// index
	mux.HandleFunc("/", controller.HandleLogin)
	// login
	mux.HandleFunc("/login", controller.HandleLogin)
	// logout
	mux.HandleFunc("/logout", controller.HandleLogout)
	enroleChecking := alice.New(controller.SessionChecker, controller.AuthEnrol)
	anyChecking := alice.New(controller.SessionChecker, controller.AuthAny)
	// work
	mux.Handle("/work", anyChecking.ThenFunc(controller.HandleWork))
	// find
	mux.Handle("/find", enroleChecking.ThenFunc(controller.Find))
	// show enrol form
	mux.Handle("/enrol/show", enroleChecking.ThenFunc(controller.ShowEnrol))
	// process enrol form
	mux.Handle("/enrol/submit", enroleChecking.ThenFunc(controller.SubmitEnrol))
	// process edit form
	mux.Handle("/edit/submit", enroleChecking.ThenFunc(controller.SubmitEdit))
	// register
	mux.HandleFunc("/register", controller.HandleRegistration)

	// konfiguriere server
	server := &http.Server{
		Addr: "127.0.0.1:8090",
		//Addr:    "0.0.0.0:8090",
		Handler: mux,
	}
	// und starte ihn
	server.ListenAndServe()
}
