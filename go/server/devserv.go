package main

import (
	"fmt"
	"github.com/geobe/gostip/go/controller"
	"github.com/geobe/gostip/go/model"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"net/http"
	"os"
	"github.com/justinas/nosurf"
	"log"
)

const pagedir = model.Base + "/pages/"
const resourcedir = model.Base + "/resources/"

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
	docbase += "/"
	// FileServer ist ein Handler, der dieses Verzeichnis bedient
	// Funktionsvariablen für alice anlegen
	files := http.StripPrefix("/pages/", http.FileServer(http.Dir(docbase + pagedir)))
	resources := http.FileServer(http.Dir(docbase + resourcedir))
	pages := http.FileServer(http.Dir(docbase + pagedir))

	requestLogging := alice.New(controller.RequestLogger)
	csrfChecking := alice.New(nosurf.NewPure)
	resultsChecking := alice.New(controller.RequestLogger, nosurf.NewPure, controller.SessionChecker, controller.AuthProjectOffice)
	enroleChecking := alice.New(controller.RequestLogger, nosurf.NewPure, controller.SessionChecker, controller.AuthEnrol)
	anyChecking := alice.New(controller.RequestLogger, nosurf.NewPure, controller.SessionChecker, controller.AuthAny)

	// Zugriff auf das Verzeichnis via Präfic /pages/
	mux.PathPrefix("/pages/").Handler(requestLogging.Then(files))
	// Zugriff auf die Resourcen-Verzeichnisse mit regular expression
	mux.PathPrefix("/{dir:(css|fonts|js|images)}/").Handler(requestLogging.Then(resources))
	// Zugriff auf *.htm[l] Dateien im /pages Verzeichnis
	mux.Handle("/{dir:\\w+\\.html?}", requestLogging.Then(pages))
	// error
	mux.HandleFunc("/err", err)
	// index
	mux.Handle("/", csrfChecking.ThenFunc(controller.HandleIndex))
	mux.Handle("/index", csrfChecking.ThenFunc(controller.HandleIndex))
	// login
	mux.Handle("/login", csrfChecking.ThenFunc(controller.HandleLogin))
	// logout
	mux.HandleFunc("/logout", controller.HandleLogout)
	// work
	mux.Handle("/work", anyChecking.ThenFunc(controller.HandleWork))
	// find
	mux.Handle("/find/applicant", enroleChecking.ThenFunc(controller.FindApplicant))
	// show results edit form
	mux.Handle("/results/show", resultsChecking.ThenFunc(controller.ShowResults))
	// submit results edit form
	mux.Handle("/results/submit", resultsChecking.ThenFunc(controller.SubmitResults))
	// show enrol form
	mux.Handle("/enrol/show", enroleChecking.ThenFunc(controller.ShowEnrol))
	// process enrol form
	mux.Handle("/enrol/submit", enroleChecking.ThenFunc(controller.SubmitEnrol))
	// show cancellation form
	mux.Handle("/cancellation/show", enroleChecking.ThenFunc(controller.ShowCancellation))
	// process cancellation form
	mux.Handle("/cancellation/submit", enroleChecking.ThenFunc(controller.SubmitCancelation))
	// process edit form
	mux.Handle("/edit/submit", enroleChecking.ThenFunc(controller.SubmitApplicantEdit))
	// register
	mux.Handle("/register", csrfChecking.ThenFunc(controller.ShowRegistration))
	// register
	mux.Handle("/register/submit", csrfChecking.ThenFunc(controller.SubmitRegistration))

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
