package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

//TLS config struct
type TLSStruct struct {
	Cert string
	Key  string
}

//
type WeiXinAppStruct struct {
	Id      string
	Secrete string
}

//ConfigStruct decide struct of config
type ConfigStruct struct {
	Sql         SqlConfig
	ProjectPath string
	TLS         TLSStruct
	WeiXinApp   WeiXinAppStruct
}

var configStrut ConfigStruct

//Sql contain sql config
//Using in database
var (
	Sql         SqlConfig
	ProjectPath string
	TLS         TLSStruct
)

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
	TLS = configStrut.TLS
}
