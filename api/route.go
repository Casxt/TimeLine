package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Casxt/TimeLine/components/line"
	"github.com/Casxt/TimeLine/components/profile"
	"github.com/Casxt/TimeLine/components/signin"
	"github.com/Casxt/TimeLine/components/signup"
)

//Route 对api进行路由
func Route(res http.ResponseWriter, req *http.Request) {
	subPath := req.URL.Path[len("/api"):] //url are like /api/(something)
	var jsonRes interface{}               //map[string]string
	var resCode int

	switch {
	//Sign Up Api
	case strings.HasPrefix(subPath, "/SignUp"):
		resCode, jsonRes = signup.SignUp(res, req)
		//Sign In Api
	case strings.HasPrefix(subPath, "/CheckAccount"):
		resCode, jsonRes = signin.CheckAccount(res, req)
	case strings.HasPrefix(subPath, "/SignIn"):
		resCode, jsonRes = signin.SignIn(res, req)
		//Line Api
	case strings.HasPrefix(subPath, "/CreateLine"):
		resCode, jsonRes = line.CreateLine(res, req)
	case strings.HasPrefix(subPath, "/AddSlice"):
		resCode, jsonRes = line.AddSlice(res, req)
	case strings.HasPrefix(subPath, "/GetSlices"):
		resCode, jsonRes = line.GetSlices(res, req)
	case strings.HasPrefix(subPath, "/GetLines"):
		resCode, jsonRes = line.GetLines(res, req)
	case strings.HasPrefix(subPath, "/GetLineInfo"):
		resCode, jsonRes = line.GetLineInfo(res, req)
	case strings.HasPrefix(subPath, "/AddUser"):
		resCode, jsonRes = line.AddUser(res, req)
		//Profile Api
	case strings.HasPrefix(subPath, "/ProfilePicture"):
		resCode, jsonRes = profile.UpdateProfilePic(res, req)
	case strings.HasPrefix(subPath, "/GetUserInfo"):
		resCode, jsonRes = profile.GetUserInfo(res, req)
	case strings.HasPrefix(subPath, "/ChangeNickName"):
		resCode, jsonRes = profile.ChangeNickName(res, req)

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
