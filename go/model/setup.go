package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sync"
	"time"
)

// the relative location of project files
const Base = "src/github.com/geobe/gostip"

// setting up viper configuration lib
func Setup(cfgfile string) {
	if cfgfile == "" {
		cfgfile = "devconfig"
	}
	viper.SetConfigName(cfgfile)
	viper.AddConfigPath(Base + "/config")    // for config in the working directory
	viper.AddConfigPath("../../config")    // for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

var db *gorm.DB
var dbsem, oblsem sync.Mutex
var oblasts []Oblast
var mailaccount, mailpw string

// make database connection available as a singleton
func Db() *gorm.DB {
	dbsem.Lock()
	defer dbsem.Unlock()
	if db == nil {
		db = ConnectDb()
	}
	return db
}

// connect to database using values from config file
func ConnectDb() *gorm.DB {
	db, err := gorm.Open(viper.GetString("db.type"), viper.GetString("db.connect"))
	if err != nil {
		fmt.Errorf("db error is %s", err)
		panic("failed to connect database with err " + fmt.Sprint(err))
	}
	return db
}

// create some test applicants in a developpment db
func InitProdDb(db *gorm.DB) *gorm.DB {

	// Migrate the schema
	db.AutoMigrate(classes...)

	// initialize if db empty
	var nOblasts int
	if db.Model(&Oblast{}).Count(&nOblasts); nOblasts < 1 {
		for _, oblast := range InitialOblasts {
			db.Create(&oblast)
		}
	}
	// retrieve oblasts from db
	var oblasts []Oblast
	db.Find(&oblasts)

	//if db.Model(&ExamReference{}).Count(&nxref); nxref < 1 {
	//	xref := ExamReference{
	//		Year:    time.Now().Year(),
	//		Results: [NQESTION]int{50, 50, 20, 30, 100, 150, 100, 100, 0, 0},
	//	}
	//	db.Create(&xref)
	//}

	users, maxresults := InitialValues()
	var aUser User
	for _, user := range users {
		db.First(&aUser, "Login = ?", user.Login)
		if aUser.Login == user.Login {
			aUser.Password = user.Password
			aUser.Role = user.Role
			db.Save(&aUser)
		} else {
			db.Save(&user)
		}
	}

	var r [NQESTION]int
	copy(r[:], maxresults)
	var xref ExamReference
	db.First(&xref, "Year = ?", time.Now().Year())
	if xref.ID != 0 {
		xref.Results = r
		xref.UpdatedBy = ""
		db.Save(&xref)
	} else {
		xref := ExamReference{
			Year:    time.Now().Year(),
			Results: r,
		}
		db.Create(&xref)
	}

	return db
}

func SetMailer(maccount, mpw string) {
	mailaccount = maccount
	mailpw = mpw
}

func CanMail() bool {
	return mailaccount != "" && mailpw != ""
}

func GetMailer() (string, string) {
	return mailaccount, mailpw
}

// make some initial users available for testing and bootstrap
func InitialValues() (users []User, maxresults []int) {
	uv := viper.Get("users").([]interface{})
	for _, v := range uv {
		switch user := v.(type) {
		case map[string]interface{}:
			nu := NewUser(
				user["login"].(string),
				user["fullname"].(string),
				user["password"].(string),
				roles(user["roles"].([]interface{}))...)
			users = append(users, *nu)
		default:
			for k1, v1 := range v.(map[string]interface{}) {
				fmt.Errorf("Error in configuration filet%s: %s\n", k1, v1)

			}
		}
	}
	mr := viper.Get("maxresults").([]interface{})
	maxresults = make([]int, len(mr))
	for i, v := range mr {
		maxresults[i] = int(v.(float64) + 0.5)
	}
	return
}

// cache oblasts in memory instead of reading every time from database
func Oblasts() []Oblast {
	oblsem.Lock()
	defer oblsem.Unlock()
	if len(oblasts) == 0 {
		Db().Order("id asc").Find(&oblasts)
	}
	return oblasts
}

// helper function to read roles from config file
func roles(r []interface{}) (res []int) {
	for _, rv := range r {
		switch rt := rv.(type) {
		case int:
			res = append(res, rt)
		case float64:
			res = append(res, int(rt))
		}
	}
	return
}
