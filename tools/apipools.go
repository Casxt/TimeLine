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

//GetLoginStateOfOperator will check wether user have login and check Operator
//and return ID and session
func GetLoginStateOfOperator(req *http.Request, SssionID, Operator string) (UserID string, Session session.IO) {
	s := session.Get(SssionID, req)
	if s == nil {
		return "", nil
	}
	ID, ok := s.Get("ID")
	if !ok {
		return "", nil
	}
	if phone, ok := s.Get("Phone"); !ok || phone != Operator {
		return "", nil
	}
	return ID, s
}

//GetLoginStateOfCookie will check wether user have login and check Cookie["Phone"]
//and return ID and session
func GetLoginStateOfCookie(req *http.Request) (UserID string, Session session.IO) {
	//Check SessionID
	c, err := req.Cookie("SessionID")
	if err != nil {
		return "", nil
	}
	s := session.Get(c.Value, req)
	if s == nil {
		return "", nil
	}
	//Check UserID
	ID, ok := s.Get("ID")
	if !ok {
		return "", nil
	}
	//Check Phone
	c, err = req.Cookie("Phone")
	if err != nil {
		return "", nil
	}
	if phone, ok := s.Get("Phone"); !ok || phone != c.Value {
		return "", nil
	}
	return ID, s
}
