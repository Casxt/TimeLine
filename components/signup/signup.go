package signup

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/page"
)

//Route decide return which Page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	switch req.Method {
	case "GET":
		status, result = Page()
	default:
		status, result = Page()
	}
	res.WriteHeader(status)
	res.Write(result)
}

//SignUp is a api interface, will signup a user
func SignUp(req *http.Request) (status int, res []byte) {
	var jsonRes map[string]string
	status = 200

	Phone := req.PostFormValue("Phone")
	Mail := req.PostFormValue("Mail")
	HashPass := req.PostFormValue("HashPass")

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
	return status, res
}

//Page will create the signup page
func Page() (status int, res []byte) {
	var err error
	res, err = page.GetPage("signup.html")
	status = 200
	if err != nil {
		status = 404
		res, _ = page.GetPage("404.html")
	}
	return status, res
}
