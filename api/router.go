package api

import (
	"log"
	"net/http"
	"strings"
)

//Route 对api进行路由
func Route(res http.ResponseWriter, req *http.Request) {

	//解析get和post字串
	if err := req.ParseForm(); err != nil {
		log.Println(err.Error())
		return
	}
	//api路由
	subPath := req.URL.Path[5:] // /api/(something)
	var jsonRes []byte
	var resCode int
	switch {
	case strings.HasPrefix(subPath, "SignUp"):
		resCode, jsonRes = SignUp(req)
	case strings.HasPrefix(subPath, "a"):
		jsonRes = []byte("asdasdsa")
	default:
		http.SetCookie(res, &http.Cookie{Name: "testcookiename2", Value: "testcookievalue", Path: "/", MaxAge: 86400})
		jsonRes = []byte("Add Cookie")
	}
	res.WriteHeader(resCode)
	res.Write(jsonRes)
}
