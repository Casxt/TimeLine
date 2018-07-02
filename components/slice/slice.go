package slice

import (
	"net/http"

	"github.com/Casxt/TimeLine/tools"
)

//AddSlice is a Api interface
//add slice to a line
func AddSlice(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	var err error

	type Data struct {
		Operator   string   `json:"Operator"`
		LineName   string   `json:"LineName"`
		Title      string   `json:"LineName"`
		Content    string   `json:"Content"`
		Gallery    []string `json:"Gallery"`
		Location   string   `json:"Location"`
		Type       string   `json:"Type"`
		Visibility string   `json:"Visibility"`
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	ID, Session := tools.GetLoginState(req)
	if Session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User  not SignIn",
		}
	}
	if Phone, ok := Session.Get("Phone"); !ok || Phone != data.Operator {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "User  not SignIn",
		}
	}

	//TODO sql
	jsonRes = map[string]string{
		"State": "Success",
		"Msg":   "slice add success",
	}
	return 200, jsonRes
}
