package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Casxt/TimeLine/components/line"
	"github.com/Casxt/TimeLine/components/signin"
	"github.com/Casxt/TimeLine/components/signup"
)

//Route 对api进行路由
func Route(res http.ResponseWriter, req *http.Request) {
	subPath := req.URL.Path[5:] //url are like /api/(something)
	var jsonRes interface{}     //map[string]string
	var resCode int

	switch {
	case strings.HasPrefix(subPath, "SignUp"):
		resCode, jsonRes = signup.SignUp(res, req)
	case strings.HasPrefix(subPath, "CheckAccount"):
		resCode, jsonRes = signin.CheckAccount(res, req)
	case strings.HasPrefix(subPath, "SignIn"):
		resCode, jsonRes = signin.SignIn(res, req)
	case strings.HasPrefix(subPath, "CreateLine"):
		resCode, jsonRes = line.CreateLine(res, req)
	case strings.HasPrefix(subPath, "AddSlice"):
		resCode, jsonRes = line.AddSlice(res, req)
	case strings.HasPrefix(subPath, "GetSlices"):
		resCode, jsonRes = line.GetSlices(res, req)
	default:
		resCode = 200
		jsonRes = map[string]string{
			"State": "Succesful",
			"Msg":   "Api is working",
		}
	}

	result, err := json.Marshal(jsonRes)

	if err != nil {
		log.Println(err)
	}

	res.WriteHeader(resCode)
	res.Write(result)
}
