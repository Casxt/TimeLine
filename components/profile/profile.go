package profile

import (
	"log"
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/static"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/profile"):]
	switch {
	default:
		status, result, _ = static.GetPage("components", "profile", "profile.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

func UpdateProfilePic(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type ReqData struct {
		SessionID string
		Operator  string
		Picture   string //Picture Hash
	}
	type ResData struct {
		State string
		Msg   string
	}
	var reqData ReqData

	if status, jsonRes = tools.GetPostJSON(req, &reqData); status != 200 {
		return status, jsonRes
	}

	UserID, session := tools.GetLoginStateOfOperator(req, reqData.SessionID, reqData.Operator)
	if session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   `User Not Sign In`,
		}
	}

	if !tools.CheckImgHash(reqData.Picture) {
		return 400, ResData{
			State: "Failed",
			Msg:   "Invalid Parameter Picture",
		}
	}
	// Check Hash Belong to User
	_, h, w, _, err := database.GetImgInfo(UserID, reqData.Picture)
	if err != nil {
		log.Println(err.Error())
		return 400, ResData{
			State: "Failed",
			Msg:   err.Error(),
		}
	}
	//Check img height=width
	if h != w {
		return 400, ResData{
			State: "Failed",
			Msg:   "Invalid Picture Size",
		}
	}
	database.UpdateProfilePic(UserID, reqData.Picture)
	return 400, ResData{
		State: "Success",
		Msg:   "Update Profile Picture Success",
	}
}
