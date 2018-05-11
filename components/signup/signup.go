package signup

import (
	"encoding/json"
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/mail"
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
func SignUp(req *http.Request) (status int, jsonRes map[string]string) {

	type SignUpData struct {
		Phone    string `json:"Phone"`
		Mail     string `json:"Mail"`
		HashPass string `json:"HashPass"`
	}
	var data SignUpData
	err := json.NewDecoder(req.Body).Decode(&data)
	if err != nil {
		status = 400
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Json bytes or parmeter not enough",
		}
		return status, jsonRes
	}

	if data.Phone == "" || data.Mail == "" {
		//参数错误
		status = 400
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Parameter",
		}
		return status, jsonRes
	}

	_, Pass, err := database.CreateUser(data.Phone, data.Mail, data.HashPass)
	if err != nil {
		//内部错误
		switch err.Error() {
		case "User Already Exist":
			status = 409
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   "User Phone or Mail Already Be Taken",
			}
		case "User Create Failde":
			status = 409
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   "Unknow Signup Error",
			}
		default:
			status = 500
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   "Unknow Signup Error",
			}
		}
		return status, jsonRes
	}

	if data.HashPass == "" {
		//计算一个初始的随机pass并Hash
		mail.SendMail(data.Mail, "TimeLine 注册验证", "<h1>TimeLine密码:"+Pass+"</h1>", nil)
	} else {
		//计算一个初始的随机pass
		mail.SendMail(data.Mail, "TimeLine 注册验证", "<h1>您已注册TimeLine</h1>", nil)
	}

	jsonRes = map[string]string{
		"State": "Success",
		"Msg":   "Create User Successful",
	}

	return 200, jsonRes
}

//Page will create the signup page
func Page() (status int, res []byte) {
	var err error
	res, err = page.GetPage("components", "signup", "signup.html")
	if err != nil {
		res, _ = page.GetPage("404.html")
		return 404, res
	}
	return 200, res
}
