package signin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/page"
	"github.com/Casxt/TimeLine/session"
	"github.com/Casxt/TimeLine/tools"
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
func CheckAccount(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	var err error

	type Data struct {
		Account string `json:"Account"`
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	if data.Account == "" {
		status = 400
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Json bytes or parmeter not enough",
		}
		return status, jsonRes
	}

	mailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegexp := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)

	var NickName, Salt, SaltPass string
	switch {
	case phoneRegexp.MatchString(data.Account):
		_, NickName, _, Salt, SaltPass, _, _, err = database.GetUserByPhone(data.Account)
	case mailRegexp.MatchString(data.Account):
		_, NickName, _, Salt, SaltPass, _, _, err = database.GetUserByMail(data.Account)
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

	session, _ := session.Auto(res, req)
	if session == nil {
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Session Create Filed",
		}
		return 500, jsonRes
	}

	rnd, _ := rand.Int(rand.Reader, big.NewInt(1<<63-1))
	session.Put("SignInVerify", rnd.String(), 300)
	session.Put("SaltPass", SaltPass, 300)
	jsonRes = map[string]string{
		"State":        "Success",
		"Msg":          "User Account Exist",
		"NickName":     NickName,
		"Salt":         Salt,
		"SignInVerify": rnd.String(),
	}
	return 200, jsonRes
}

//SignIn Will Auth User in a indirect way
func SignIn(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		Encrypted string `json:"Encrypted"`
		IV        string `json:"IV"`
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}
	session := session.GetByCookie(req)
	if session == nil {
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg": `Session not found, you should have SessionID first,
			 call CheckAccount before you call this api`,
		}
		return 400, jsonRes
	}

	var SaltPass, SignInVerify string
	var ok bool
	if SaltPass, ok = session.Get("SaltPass"); !ok {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "SaltPass of Session not found",
		}
	}
	session.Delete(SaltPass)
	if SignInVerify, ok = session.Get("SignInVerify"); !ok {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "SignInVerify of Session not found",
		}
	}
	session.Delete(SignInVerify)
	var key, ciphertext, nonce []byte
	var err error

	if key, err = hex.DecodeString(SaltPass); err != nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invaild Parameter SaltPass",
		}
	}

	if ciphertext, err = hex.DecodeString(data.Encrypted); err != nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invaild Parameter Encrypted",
		}
	}

	if nonce, err = hex.DecodeString(data.IV); err != nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invaild Parameter IV",
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Decrypted failed",
		}
	}
	if string(plaintext) != SignInVerify {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Compare failed",
		}
	}
	/*
		aesgcm, err = cipher.NewGCM(block)
		if err != nil {
			panic(err.Error())
		}
		temp := []byte(SignInVerify)
		serverCiphertext := aesgcm.Seal(nil, nonce, temp, nil)

		if string(serverCiphertext) != string(ciphertext) {
			return 400, map[string]string{
				"State": "Failde",
				"Msg":   "Compare failed",
			}
		}
	*/
	return 200, map[string]string{
		"State": "Success",
		"Msg":   "Sign In Success",
	}
}
