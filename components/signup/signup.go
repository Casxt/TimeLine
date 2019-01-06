package signup

import (
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/mail"
	"github.com/Casxt/TimeLine/static"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide return which static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/signup"):]

	switch {
	//case strings.HasSuffix(strings.ToLower(subPath), "signup.js"):
	//	result, status, _ = static.GetFile("components", "signup", "signup.js")
	//	res.Header().Add("Content-Type", "application/x-javascript")
	default:
		status, result, _ = static.GetPage("components", "signup", "signup.html")

	}
	res.WriteHeader(status)
	res.Write(result)
}

//SignUp is a api interface, will signup a user
func SignUp(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {

	type Data struct {
		Phone    string
		Mail     string
		HashPass string
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

	_, Pass, err := database.CreateUser(data.Phone, data.Mail, data.HashPass, nil, nil)
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
		if mail.SendMail(data.Mail, "TimeLine 注册验证", "<h1>TimeLine密码:"+Pass+"</h1>", nil) != nil {
			return 409, map[string]string{
				"State": "Failed",
				"Msg":   "SendMail unKnow failed",
			}
		}
	} else {
		//计算一个初始的随机pass
		if mail.SendMail(data.Mail, "TimeLine 注册验证", "<h1>您已注册TimeLine</h1>", nil) != nil {
			return 409, map[string]string{
				"State": "Failed",
				"Msg":   "SendMail unKnow failed",
			}
		}
	}

	return 200, map[string]string{
		"State": "Success",
		"Msg":   "Create User Successful",
	}
}

func WeiXinSignUp(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {

	type Data struct {
		UserCode string
		Phone    string
		Mail     string
		HashPass string
	}
	var data Data

	var (
		OpenID  string
		UnionID *string
	)
	if OpenID, _, UnionID, status, jsonRes = tools.GetWeiXinUser(data.UserCode); status != 200 {
		return status, jsonRes
	}

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	if data.Phone == "" || data.Mail == "" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Parameter",
		}
	}

	_, Pass, err := database.CreateUser(data.Phone, data.Mail, data.HashPass, &OpenID, UnionID)
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
		if mail.SendMail(data.Mail, "TimeLine 注册验证", "<h1>TimeLine密码:"+Pass+"</h1>", nil) != nil {
			return 409, map[string]string{
				"State": "Failed",
				"Msg":   "SendMail unKnow failed",
			}
		}
	} else {
		//计算一个初始的随机pass
		if mail.SendMail(data.Mail, "TimeLine 注册验证", "<h1>您已注册TimeLine</h1>", nil) != nil {
			return 409, map[string]string{
				"State": "Failed",
				"Msg":   "SendMail unKnow failed",
			}
		}
	}

	return 200, map[string]string{
		"State": "Success",
		"Msg":   "Create User Successful",
	}
}

func BindingWeiXinToUser(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		UserCode string
		Phone    string
		Mail     string
	}
	var data Data

	var (
		OpenID  string
		UnionID *string
	)
	if OpenID, _, UnionID, status, jsonRes = tools.GetWeiXinUser(data.UserCode); status != 200 {
		return status, jsonRes
	}

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	if data.Phone == "" || data.Mail == "" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Parameter",
		}
	}

	if DBErr := database.UpdateWeiXin(OpenID, data.Phone, data.Mail, UnionID); DBErr != nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   DBErr.Error(),
		}
	}
	return 200, map[string]string{
		"State": "Success",
		"Msg":   "Binding User Successful",
	}
}
