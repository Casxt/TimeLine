package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//SqlConfig contain sql config
type SqlConfig struct {
	Host string
	User string
	Pass string
	Port string
	Base string
}

//ConfigStrut decide struct of config
type ConfigStrut struct {
	Sql SqlConfig
}

var configStrut ConfigStrut

//Sql contain sql config
//Using in database
var Sql *SqlConfig

//Load Config file
func Load(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	log.Println("Read Config \r\n", string(data))
	err = json.Unmarshal(data, &configStrut)
	if err != nil {
		panic(err)
	}
	Sql = &configStrut.Sql
}
