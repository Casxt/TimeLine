package signin

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/session"
	"github.com/Casxt/TimeLine/static"
	"github.com/Casxt/TimeLine/tools"
	"html"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
)

//Route Return The static to Show
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/signin"):]
	switch {
	//case strings.HasPrefix(strings.ToLower(subPath), "signin.js"):
	//	result, status, _ = static.GetFile("components", "signin", "signup.js")
	default:
		status, result, _ = static.GetPage("components", "signin", "signin.html")

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

	var ID, Phone, NickName, Salt, SaltPass string
	switch {
	case phoneRegexp.MatchString(data.Account):
		ID, _, NickName, _, Salt, SaltPass, _, _, err = database.GetUserByPhone(data.Account, nil)
		Phone = data.Account
	case mailRegexp.MatchString(data.Account):
		ID, Phone, NickName, _, Salt, SaltPass, _, _, err = database.GetUserByMail(data.Account)
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
	//Save some info of account, if login failed, it will be delete
	session.Put("ID", ID, 300)
	session.Put("NickName", NickName, 300)
	session.Put("Phone", Phone, 300)
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
		Encrypted string //加密后的字符串,16进制字符串
		IV        string //加密后的IV,16禁止字符串
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}
	sess := session.GetByCookie(req)
	if sess == nil {
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg": `Session not found, you should have SessionID first,
			 call CheckAccount before you call this api`,
		}
		return 400, jsonRes
	}

	var SaltPass, SignInVerify string
	var ok bool
	if SaltPass, ok = sess.Get("SaltPass"); !ok {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "SaltPass of Session not found",
		}
	}
	sess.Delete(SaltPass)
	if SignInVerify, ok = sess.Get("SignInVerify"); !ok {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "SignInVerify of Session not found",
		}
	}
	sess.Delete(SignInVerify)
	var key, ciphertext, nonce []byte
	var err error

	if key, err = hex.DecodeString(SaltPass); err != nil {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invaild Parameter SaltPass",
		}
	}

	if ciphertext, err = hex.DecodeString(data.Encrypted); err != nil {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invaild Parameter Encrypted",
		}
	}

	if nonce, err = hex.DecodeString(data.IV); err != nil {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
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
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Decrypted failed",
		}
	}
	if string(plaintext) != SignInVerify {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Compare failed",
		}
	}
	ID, ok := sess.Get("ID")
	if !ok {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Timeout",
		}
	}
	NickName, ok := sess.Get("NickName")
	if !ok {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Timeout",
		}
	}
	Phone, ok := sess.Get("Phone")
	if !ok {
		sess.Delete("ID")
		sess.Delete("NickName")
		sess.Delete("Phone")
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Timeout",
		}
	}
	sess.Put("ID", ID, 60*60*24*30)
	sess.Put("NickName", NickName, 60*60*24*30)
	sess.Put("Phone", Phone, 60*60*24*30)
	//Escape, Phone is pure number and no need to escape
	http.SetCookie(res, &http.Cookie{Name: "Phone", Value: Phone, Path: "/", MaxAge: 30 * 86400})
	http.SetCookie(res, &http.Cookie{Name: "NickName", Value: html.EscapeString(url.QueryEscape(NickName)), Path: "/", MaxAge: 30 * 86400})

	return 200, map[string]string{
		"State":     "Success",
		"Msg":       "Sign In Success",
		"SessionID": sess.ID(),
		"Phone":     Phone,
		"NickName":  NickName,
	}
}

//WeiXinSignIn 微信小程序登录接口
func WeiXinSignIn(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		UserCode string //用户身份码
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	var OpenID, SessionKey string
	if OpenID, SessionKey, _, status, jsonRes = tools.GetWeiXinUser(data.UserCode); status != 200 {
		return status, jsonRes
	}

	ID, _, Phone, _, NickName, _, _, _, _, _, DBErr := database.GetUserByOpenID(OpenID)
	if DBErr != nil {
		//"User Not Exist"
		return 400, map[string]string{
			"State": "Failed",
			"Msg":   DBErr.Error(),
		}
	}

	sess := session.New(req)
	if sess == nil {
		return 500, map[string]string{
			"State": "Failed",
			"Msg":   "Unknown error",
		}
	}
	sess.Put("ID", ID, 60*60*24*30)
	sess.Put("NickName", NickName, 60*60*24*30)
	sess.Put("Phone", Phone, 60*60*24*30)
	sess.Put("WeiXinSession", SessionKey, 60*60*24*30)
	return 200, map[string]string{
		"State":     "Success",
		"Msg":       "Sign In Success",
		"SessionID": sess.ID(),
		"Phone":     Phone,
		"NickName":  NickName,
	}
}
