package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"sync"
)

func Setup(cfgfile string) {
	if cfgfile == "" {
		cfgfile = "devconfig.json"
	}
	viper.SetConfigFile(cfgfile)
	viper.AddConfigPath(".")    // for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

var db *gorm.DB
var dbsem, oblsem sync.Mutex
var oblasts []Oblast

func Db() *gorm.DB {
	dbsem.Lock()
	defer dbsem.Unlock()
	if db == nil {
		db = ConnectDb()
	}
	return db
}

func ConnectDb() *gorm.DB {
	db, err := gorm.Open(viper.GetString("db.type"), viper.GetString("db.connect"))
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

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

func Oblasts() []Oblast {
	oblsem.Lock()
	defer oblsem.Unlock()
	if len(oblasts) == 0 {
		Db().Find(&oblasts)
	}
	return oblasts
}

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
