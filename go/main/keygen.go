package main

import (
	"github.com/gorilla/securecookie"
	"fmt"
	"github.com/geobe/gostip/go/model"
	"github.com/spf13/viper"
)

func main() {
	k := securecookie.GenerateRandomKey(32)
	fmt.Printf("%v\n",k)
	model.Setup("")

	key := viper.Get("csrfkey").([]interface{})
	//len := len(key)
	var bkey [32]byte
	for i, v := range key {
		bkey[i] = byte(v.(float64))
	}
	sec  := viper.GetBool("csrfsecure")
	fmt.Printf("secure: %t, key: %v\n", sec, key)
	fmt.Printf("secure: %t, key: %v\n", sec, bkey)
}
