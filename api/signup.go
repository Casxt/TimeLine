package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MapleFadeAway/timeline/database"
)

//SignUpInfo is the info needed when register
type SignUpInfo struct {
	Phone    string `json:"Phone"`
	Mail     string `json:"Mail"`
	HashPass string `json:"HashPass"`
}

//SignUp will signup a user
func SignUp(req *http.Request) (status int, res []byte) {
	var jsonRes map[string]string
	status = 200

	Phone := req.FormValue("Phone")
	Mail := req.FormValue("Mail")
	HashPass := req.FormValue("HashPass")

	if Phone == "" || Mail == "" || HashPass == "" {
		//参数错误
		status = 400
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Parameter",
		}
	}

	if err := database.CreateUser(Phone, Mail, HashPass); err != nil {
		//内部错误
		status = 500
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   err.Error(),
		}
	} else {
		jsonRes = map[string]string{
			"State": "Succesful",
			"Msg":   "Create User Successful",
		}
	}

	res, err := json.Marshal(jsonRes)

	if err != nil {
		log.Fatalln(err)
	}
	return 0, res
}
