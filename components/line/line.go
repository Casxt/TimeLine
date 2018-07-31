package line

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/static"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	subPath := req.URL.Path[len("/line"):]
	switch {
	//match /linecreate*
	case strings.HasPrefix(strings.ToLower(subPath), "create"):
		status, result, _ = static.GetPage("components", "line", "createLine.html")
	//match /lineedit*
	case strings.HasPrefix(strings.ToLower(subPath), "edit"):
		status, result, _ = static.GetPage("components", "line", "editLine.html")
	default:
		status, result, _ = static.GetPage("components", "line", "line.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

//GetLines Get Lines of User
func GetLines(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type ReqData struct {
		Operator  string
		SessionID string
	}
	type ResData struct {
		State string
		Msg   string
		Lines []string
	}
	var reqData ReqData

	if status, jsonRes = tools.GetPostJSON(req, &reqData); status != 200 {
		return status, jsonRes
	}

	UserID, _ := tools.GetLoginStateOfOperator(req, reqData.SessionID, reqData.Operator)
	if UserID == "" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User Not Sign In",
		}
	}

	Lines, err := database.GetLines(UserID)
	if err != nil {
		log.Println(err.Error())
		return 500, map[string]string{
			"State":  "Failde",
			"Msg":    "database.GetUserLines Failed",
			"Detial": err.Error(),
		}
	}

	return 200, ResData{
		State: "Success",
		Msg:   "Lines Get Successful",
		Lines: Lines,
	}
}

//GetLineInfo return LineInfo include some sum info
func GetLineInfo(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type ReqData struct {
		Operator  string
		SessionID string
		LineName  string
	}
	type ResData struct {
		State      string
		Msg        string
		Name       string
		LatestImg  string
		Users      []string
		SliceNum   int
		ImageNum   int
		CreateTime time.Time
		LatestTime time.Time
	}

	var (
		reqData ReqData
		resData ResData
		err     error
	)

	if status, jsonRes = tools.GetPostJSON(req, &reqData); status != 200 {
		return status, jsonRes
	}

	UserID, _ := tools.GetLoginStateOfOperator(req, reqData.SessionID, reqData.Operator)
	if UserID == "" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User Not Sign In",
		}
	}
	_, resData.Name, resData.LatestImg, resData.Users, resData.SliceNum, resData.ImageNum,
		resData.CreateTime, resData.LatestTime, err = database.GetLineDetail(reqData.LineName, nil)
	if err != nil {
		return 500, map[string]string{
			"State":  "Failde",
			"Msg":    "Get Line Info Failed",
			"Detail": err.Error(),
		}
	}
	resData.State = "Success"
	resData.Msg = "Get Line Info Success"
	return 200, resData
}

//CreateLine will create a new line with specific name
func CreateLine(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		Operator  string
		LineName  string
		SessionID string
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, session := tools.GetLoginStateOfOperator(req, data.SessionID, data.Operator)
	if session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   `User Not Sign In`,
		}
	}

	if len(data.LineName) < 4 {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Line Name too short",
		}
	}

	Lines, err := database.GetLines(UserID)
	if err != nil {
		log.Println(err.Error())
		return 400, map[string]string{
			"State":  "Failde",
			"Msg":    "Unknow Error Happend",
			"Detail": err.Error(),
		}
	}

	//Limit Line num of User
	if len(Lines) > 3 {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User Have too Much Lines",
		}
	}

	err = database.CreateLine(strings.ToLower(data.LineName), UserID)
	if err != nil {
		switch err.Error() {
		case "Line Already Exist":
			return 400, map[string]string{
				"State": "Failde",
				"Msg":   "Line Name Already be Used",
			}
		default:
			return 400, map[string]string{
				"State": "Failde",
				"Msg":   "Line Name Already be Used",
			}
		}
	}
	return 200, map[string]string{
		"State": "Success",
		"Msg":   "Line create success",
	}

}

//AddUser will will add user to specific line
func AddUser(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type ReqData struct {
		Operator  string
		SessionID string
		LineName  string
		NickName  string
		UserPhone string
	}
	type ResData struct {
		State  string
		Msg    string
		Detail string
	}

	var reqData ReqData

	if status, jsonRes = tools.GetPostJSON(req, &reqData); status != 200 {
		return status, jsonRes
	}

	UserID, _ := tools.GetLoginStateOfOperator(req, reqData.SessionID, reqData.Operator)
	if UserID == "" {
		return 400, ResData{
			State: "Failde",
			Msg:   "User Not Sign In",
		}
	}
	//Check if user belong to line
	var contain bool
	lines, DBErr := database.GetLines(UserID)
	if DBErr != nil {
		return 500, ResData{
			State: "Failde",
			Msg:   "GetLines failed",
		}
	}

	for _, line := range lines {
		if line == reqData.LineName {
			contain = true
			break
		}
	}
	if !contain {
		return 400, ResData{
			State: "Failde",
			Msg:   "User Dose Not Belong to this Line",
		}
	}

	//Check User Phone and NickName match
	UserID, _, nickName, _, _, _, _, _, DBErr := database.GetUserByPhone(reqData.UserPhone, nil)
	if DBErr != nil {
		return 500, ResData{
			State:  "Failde",
			Msg:    "Get GetUserByPhone Failed",
			Detail: DBErr.Error(),
		}
	}

	if reqData.NickName != nickName {
		return 400, ResData{
			State: "Failde",
			Msg:   "Invilde Parameter",
		}
	}

	if DBErr := database.AddUser(reqData.LineName, UserID, nil); DBErr != nil {
		return 500, ResData{
			State:  "Failde",
			Msg:    "AddUser Failed",
			Detail: DBErr.Error(),
		}
	}

	return 200, ResData{
		State: "Success",
		Msg:   "User Add Success",
	}
}
