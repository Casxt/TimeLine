package line

import (
	"bytes"
	"log"
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/tools"
)

func GetSlices(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type Data struct {
		SessionID string
		Operator  string
		LineName  string
		PageNum   int //1-n
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

	Slices, err := database.GetSlices(data.LineName, UserID, data.PageNum)
	if err != nil {
		log.Println(err.Error())
		return 500, map[string]string{
			"State":  "Failde",
			"Msg":    "GetSlices Database error",
			"Detial": err.Error(),
		}
	}
	return 200, SlicesInfo{
		State:  "Failde",
		Msg:    "User  not SignIn",
		Slices: Slices,
	}
}

//AddSlice is a Api interface
//add slice to a line
func AddSlice(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	//var err error
	//Title      string   `json:"LineName"`
	type Data struct {
		Operator   string
		LineName   string
		Content    string
		Gallery    []string
		Type       string
		Visibility string
		Longitude  string //精度
		Latitude   string //纬度
		Time       string
		SessionID  string
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, Session := tools.GetLoginStateOfOperator(req, data.SessionID, data.Operator)
	if Session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User  not SignIn",
		}
	}

	//TODO: Check Type
	//TODO: Check Gallery
	//TODO: Check Content Visibility Longitude Latitude LineName

	var galleryString string
	//Gallery into hash1,hahs2,...hashn, format
	imgNum := len(data.Gallery)
	if imgNum > 0 {
		//64*imgNum+imgNum
		buff := bytes.NewBuffer(make([]byte, 65*imgNum))
		buff.Reset()
		for _, Hash := range data.Gallery {
			buff.WriteString(Hash)
			buff.Write([]byte(","))
		}
		galleryString = string(buff.Bytes()[0 : buff.Len()-1])
	} else {
		galleryString = ""
	}

	Location := data.Longitude + "," + data.Latitude
	//TODO: Check How Many Slice User have create today
	err := database.CreateSlice(data.LineName, UserID, data.Content, galleryString, data.Type, data.Visibility, Location, data.Time)
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
