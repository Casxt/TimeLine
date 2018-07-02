package tools

import (
	"encoding/json"
	"net/http"

	"github.com/Casxt/TimeLine/session"
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

//GetLoginState will check wether user have login
//and return ID and session
func GetLoginState(req *http.Request) (UserID string, Session session.IO) {
	c, err := req.Cookie("SessionID")
	if err != nil {
		return "", nil
	}
	s := session.Get(c.Value, req)
	if s == nil {
		return "", nil
	}
	ID, ok := s.Get("ID")
	if !ok {
		return "", nil
	}
	return ID, s
}
