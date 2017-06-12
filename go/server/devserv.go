package main

import (
	"github.com/geobe/gostip/go/controller"
	"github.com/geobe/gostip/go/model"
	"net/http"
	"log"
	"golang.org/x/crypto/acme/autocert"
	"crypto/tls"
	"flag"
)

// Server Ports, zu denen  Ports 80 und 443
// vom Internet Router (z.B. FritzBox) mit Port Forwarding weitergeleitet wird
const httpport = ":8070"
const tlsport = ":8443"

func main() {
	// read command line parameters
	account := flag.String("mail", "", "a mail account")
	mailpw := flag.String("mailpw", "", "password of the mail account")
	cfgfile := flag.String("cfgfile", "", "name of config file")
	flag.Parse()
	// setup mailer info
	model.SetMailer(*account, *mailpw)
	// prepare database
	model.Setup(*cfgfile)
	db := model.Db()
	defer db.Close()
//	model.ClearTestDb(db)
	model.InitProdDb(db)

	// mux verwaltet die Routen
	mux := controller.SetRouting()

	// die zugelassenen host namen
	allowedHosts := []string{"dkfai.spdns.org", "geobe.spdns.org"}

	// der Verwalter der LetsEncrypt Zertifikate
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(allowedHosts...), //your domain here
		Email: 	    "georg.beier@fh-zwickau.de",
		Cache:      autocert.DirCache("certs"), //folder for storing certificates
	}



	// konfiguriere server
	server := &http.Server{
		Addr:    "0.0.0.0" + tlsport,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
		Handler: mux,
	}

	// switching redirect handler
	handlerSwitch := &controller.HandlerSwitch{
		Mux:    mux,
		Redirect: http.HandlerFunc(controller.RedirectHTTP),
		AllowedHosts: allowedHosts,
	}

	// konfiguriere redirect server
	redirectserver := &http.Server{
		Addr:    "0.0.0.0" + httpport,
		Handler: handlerSwitch, //http.HandlerFunc(RedirectHTTP),
	}
	// starte den redirect server auf HTTP
	go redirectserver.ListenAndServe()

	// und starte den primären server auf HTTPS
	log.Printf("server starting\n")
	server.ListenAndServeTLS("", "")
}
