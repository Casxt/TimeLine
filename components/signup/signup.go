package signup

import (
	"net/http"
	"strings"

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

	Phone := req.PostFormValue("Phone")
	Mail := req.PostFormValue("Mail")
	HashPass := req.PostFormValue("HashPass")

	if Phone == "" || Mail == "" {
		//参数错误
		status = 400
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Parameter",
		}
	}

	_, Pass, err := database.CreateUser(Phone, Mail, HashPass)
	if err != nil {
		//内部错误
		if strings.HasPrefix(err.Error(), "Error 1062: Duplicate entry") {
			status = 409
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   "User Name Already Be Taken",
			}
		} else {
			status = 500
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   err.Error(),
			}
		}

		return status, jsonRes
	}

	if HashPass == "" {
		//计算一个初始的随机pass并Hash
		mail.SendMail(Mail, "TimeLine 注册验证", "<h1>TimeLine密码:"+Pass+"</h1>", nil)
	} else {
		//计算一个初始的随机pass
		mail.SendMail(Mail, "TimeLine 注册验证", "<h1>您已注册TimeLine</h1>", nil)
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
