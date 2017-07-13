package main

import (
	"github.com/geobe/gostip/go/model"
	"fmt"
	"flag"
)

func main()  {
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

	var database []*model.ApplicantData = make([]*model.ApplicantData,0)

	//db.Where("updater > ?", 0).Find(&database)

	db.Where("updater > ?", 0).Find(&database)

	for _, app :=range database {
		fmt.Println(app.UpdatedBy,app.LastName,app.FirstName,app.UpdatedAt)
	}
}
