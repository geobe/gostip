package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sync"
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
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		// Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

var db *gorm.DB
var dbsem, oblsem sync.Mutex
var oblasts []Oblast

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
		panic("failed to connect database")
	}
	return db
}

// make some initial users available for testing and bootstrap
func InitialUsers() (users []User) {
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
