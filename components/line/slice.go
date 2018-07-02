package line

import (
	"bytes"
	"net/http"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/tools"
)

//AddSlice is a Api interface
//add slice to a line
func AddSlice(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	//var err error
	//Title      string   `json:"LineName"`
	type Data struct {
		Operator   string   `json:"Operator"`
		LineName   string   `json:"LineName"`
		Content    string   `json:"Content"`
		Gallery    []string `json:"Gallery"`
		Type       string   `json:"Type"`
		Visibility string   `json:"Visibility"`
		Longitude  string   `json:"Longitude"` //精度
		Latitude   string   `json:"Latitude"`  //纬度
		Time       string   `json:"Time"`
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, Session := tools.GetLoginState(req, data.Operator)
	if Session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User  not SignIn",
		}
	}

	//TODO: Check Type
	//TODO: Check Gallery
	//TODO: Check Content Visibility Longitude Latitude LineName

	//Gallery into hash1,hahs2,...hashn, format
	imgNum := len(data.Gallery)
	//64*imgNum+imgNum
	buff := bytes.NewBuffer(make([]byte, 65*imgNum))
	for _, Hash := range data.Gallery {
		buff.Write([]byte(Hash))
		buff.Write([]byte(","))
	}
	Location := data.Longitude + "," + data.Latitude
	database.CreateSlice(data.LineName, UserID, data.Content, string(buff.Bytes()[0:buff.Len()]), data.Type, data.Visibility, Location, data.Time)
	//TODO sql
	jsonRes = map[string]string{
		"State": "Success",
		"Msg":   "slice add success",
	}
	return 200, jsonRes
}
