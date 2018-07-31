package profile

import (
	"log"
	"net/http"
	"time"

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
		return 400, ResData{
			State: "Failed",
			Msg:   "User Not Sign In",
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

func GetUserInfo(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type ReqData struct {
		SessionID string
		Operator  string
	}
	type ResData struct {
		State      string
		Msg        string
		NickName   string    //昵称
		Phone      string    //手机
		Mail       string    //邮箱
		Gender     string    //性别
		ProfilePic string    //头像
		SignInTime time.Time //注册时间
	}
	var (
		reqData ReqData
		err     error
	)

	if status, jsonRes = tools.GetPostJSON(req, &reqData); status != 200 {
		return status, jsonRes
	}

	_, session := tools.GetLoginStateOfOperator(req, reqData.SessionID, reqData.Operator)
	if session == nil {
		return 400, ResData{
			State: "Failed",
			Msg:   "User Not Sign In",
		}
	}
	resData := ResData{
		State: "Success",
		Msg:   "Get User Info Success",
		Phone: reqData.Operator,
	}
	_, resData.Mail, resData.NickName, resData.Gender, _, _, resData.ProfilePic, resData.SignInTime, err =
		database.GetUserByPhone(reqData.Operator, nil)
	if err != nil {
		log.Println(err.Error())
		return 400, ResData{
			State: "Failed",
			Msg:   err.Error(),
		}
	}

	return 200, resData
}

func ChangeNickName(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type ReqData struct {
		SessionID string
		Operator  string
		NewName   string //新昵称
	}
	type ResData struct {
		State string
		Msg   string
	}
	var reqData ReqData
	if status, jsonRes = tools.GetPostJSON(req, &reqData); status != 200 {
		return status, jsonRes
	}

	UserID, _ := tools.GetLoginStateOfOperator(req, reqData.SessionID, reqData.Operator)
	if UserID == "" {
		return 400, ResData{
			State: "Failed",
			Msg:   "User Not Sign In",
		}
	}
	if !tools.ChecNickName(reqData.NewName) {
		return 400, ResData{
			State: "Failed",
			Msg:   "Invalid Name",
		}
	}
	if err := database.UpdateNickName(UserID, reqData.NewName, nil); err != nil {
		log.Println("ChangeNickName", err.Error())
		return 500, ResData{
			State: "Failed",
			Msg:   "UpdateNickName Failed",
		}
	}
	return 200, ResData{
		State: "Success",
		Msg:   "UpdateNickName Success",
	}
}
