package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//ConfigStrut decide struct of config
type ConfigStrut struct {
	Sql         SqlConfig
	ProjectPath string
}

var configStrut ConfigStrut

//Sql contain sql config
//Using in database
var Sql SqlConfig
var ProjectPath string

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
	Sql = configStrut.Sql
	ProjectPath = configStrut.ProjectPath
}
