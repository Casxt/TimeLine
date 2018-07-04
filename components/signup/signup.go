package signup

import (
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/mail"
	"github.com/Casxt/TimeLine/page"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide return which Page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/signup"):]

	switch {
	//case strings.HasSuffix(strings.ToLower(subPath), "signup.js"):
	//	result, status, _ = page.GetFile("components", "signup", "signup.js")
	//	res.Header().Add("Content-Type", "application/x-javascript")
	default:
		status, result, _ = page.GetPage("components", "signup", "signup.html")

	}
	res.WriteHeader(status)
	res.Write(result)
}

//SignUp is a api interface, will signup a user
func SignUp(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {

	type Data struct {
		Phone    string `json:"Phone"`
		Mail     string `json:"Mail"`
		HashPass string `json:"HashPass"`
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	if data.Phone == "" || data.Mail == "" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Parameter",
		}
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
