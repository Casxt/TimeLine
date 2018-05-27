package signin

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/page"
	"github.com/Casxt/TimeLine/session"
)

//Route Return The Page to Show
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	subPath := req.URL.Path[7:]
	switch {
	case strings.HasPrefix(strings.ToLower(subPath), "signin.js"):
		result, status, _ = page.GetFile("components", "signin", "signup.js")
	default:
		result, status, _ = page.GetPage("components", "signin", "signin.html")

	}
	res.WriteHeader(status)
	res.Write(result)
}

//CheckAccount is a Api interface
//the first step of sign in is check the account
func CheckAccount(req *http.Request) (status int, jsonRes map[string]string) {
	var Salt string
	var err error
	type Data struct {
		Account string `json:"Account"`
	}
	var data Data
	err = json.NewDecoder(req.Body).Decode(&data)
	if err != nil || data.Account == "" {
		status = 400
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Json bytes or parmeter not enough",
		}
		return status, jsonRes
	}

	mailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegexp := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	var NickName string
	switch {
	case phoneRegexp.MatchString(data.Account):
		_, NickName, _, Salt, _, _, _, err = database.GetUserByPhone(data.Account)
	case mailRegexp.MatchString(data.Account):
		_, NickName, _, Salt, _, _, _, err = database.GetUserByMail(data.Account)
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

	s := session.New(req)
	if s == nil {
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Session Create Filed",
		}
		return 500, jsonRes
	}

	rnd, _ := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	s.Put("SignInVerify", rnd.String())
	jsonRes = map[string]string{
		"State":        "Success",
		"Msg":          "User Account Exist",
		"NickName":     NickName,
		"Salt":         Salt,
		"SignInVerify": rnd.String(),
	}
	return 200, jsonRes
}
