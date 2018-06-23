package tools

import (
	"encoding/json"
	"net/http"
)

//GetPostJSON get json in req, by using data.
//return 400 and errinfo if err
//ruturn 200 if success
func GetPostJSON(req *http.Request, data interface{}) (status int, jsonRes map[string]string) {
	if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
		jsonRes = map[string]string{
			"State": "Failde",
			"Msg":   "Invilde Json bytes or parmeter not enough",
		}
		return 400, jsonRes
	}
	return 200, nil
}
