package signin

import (
	"net/http"
	"regexp"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/page"
)

//Route Return The Page to Show
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

//CheckAccount is a Api interface
//the first step of sign in is check the account
func CheckAccount(req *http.Request) (status int, jsonRes map[string]string) {
	Account := req.PostFormValue("Account")

	mailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegexp := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	var Salt string
	var err error
	switch {
	case mailRegexp.MatchString(Account):
		_, _, _, Salt, _, _, _, err = database.GetUserByPhone(Account)
	case phoneRegexp.MatchString(Account):
		_, _, _, Salt, _, _, _, err = database.GetUserByMail(Account)
	default:
		//不匹配邮箱或电话
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Account",
		}
		return 400, jsonRes
	}

	//判断数据库错误
	if err != nil {
		switch err.Error() {
		case "User Not Exist":
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   err.Error(),
			}
			return 400, jsonRes
		default:
			jsonRes = map[string]string{
				"State": "Failde",
				"Msg":   err.Error(),
			}
			return 500, jsonRes
		}

	}

	jsonRes = map[string]string{
		"State": "Success",
		"Msg":   "User Account Exist",
		"Salt":  Salt,
	}
	return 200, jsonRes
}

//Page will create the signin page
func Page() (status int, res []byte) {
	var err error
	res, err = page.GetPage("components", "signin", "signin.html")
	if err != nil {
		res, _ = page.GetPage("404.html")
		return 404, res
	}
	return 200, res
}
