package line

import (
	"log"
	"net/http"
	"strings"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/tools"
)

//GetSlices GetSlices
func GetSlices(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type Data struct {
		SessionID string
		Operator  string
		LineName  string //Line名
		PageNum   int    //页数，每页20条
	}
	type SlicesInfo struct {
		State  string
		Msg    string
		Slices []database.SliceInfo
	}
	var data Data
	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, Session := tools.GetLoginStateOfOperator(req, data.SessionID, data.Operator)
	if Session == nil {
		//Do Nothing
	}
	//Check PageNum
	if data.PageNum < 1 {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "invalid PageNum",
		}
	}
	//TODO: Check LineName
	Slices, err := database.GetSlices(strings.ToLower(data.LineName), UserID, 1)
	if err != nil {
		log.Println(err.Error())
		return 500, map[string]string{
			"State":  "Failde",
			"Msg":    "GetSlices Database error",
			"Detial": err.Error(),
		}
	}
	return 200, SlicesInfo{
		State:  "Success",
		Msg:    "Slices Get!",
		Slices: Slices,
	}
}

//AddSlice is a Api interface
//add slice to a line
func AddSlice(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	//var err error
	//Title      string   `json:"LineName"`
	type Data struct {
		Operator   string   //用户手机
		LineName   string   //要添加到的Line名字
		Content    string   //Slice的文字内容
		Gallery    []string //Slice的附带图片
		Type       string   //种类，回忆还是纪念日
		Visibility string   //查看权限，自己可见，同Line可见，公开
		Longitude  string   //精度
		Latitude   string   //纬度
		Time       string   //记录时间
		SessionID  string   //SessionID
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, Session := tools.GetLoginStateOfOperator(req, data.SessionID, data.Operator)
	if Session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User not SignIn",
		}
	}

	//TODO: Check Type

	//Check Gallery
	for _, imgHash := range data.Gallery {
		if !tools.CheckImgHash(imgHash) {
			return 400, map[string]string{
				"State": "Failde",
				"Msg":   "Invalid Image Hash",
			}
		}
	}

	//TODO: Check imgHash exist

	//TODO: Check Content Longitude Latitude LineName
	if data.Visibility != "Protect" && data.Visibility != "Public" && data.Visibility != "Private" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Invalid Visibility",
		}
	}

	Location := data.Longitude + "," + data.Latitude
	//TODO: Check How Many Slice User have create today
	err := database.CreateSlice(data.LineName, UserID, data.Gallery, data.Content, data.Type, data.Visibility, Location, data.Time)
	if err != nil {
		return 200, map[string]string{
			"State": "Success",
			"Msg":   err.Error(),
		}
	}
	return 200, map[string]string{
		"State": "Success",
		"Msg":   "slice add success",
	}
}
