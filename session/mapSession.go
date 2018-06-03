package session

import (
	"net/http"
	"time"
)

var (
	sessionMap map[string]SessionIO
)

func Open() {
	sessionMap = make(map[string]SessionIO)
}

func New(req *http.Request) SessionIO {
	var session SessionIO
	var sessionID string
	sessionID = newID()
	count := 0

	for sessionMap[sessionID] != nil && count < 3 {
		sessionID = newID()
		count++
	}

	if count == 3 {
		panic("too many key duplicate")
	}

	session.Init(sessionID)
	session.ExpireTime(time.Duration(time.Hour * 24 * 30))
	session.Put("RemoteAddr", req.RemoteAddr)

	sessionMap[sessionID] = session
	return session
}

func Get(sessionID string, req *http.Request) SessionIO {

	if session, ok := sessionMap[sessionID]; ok {
		return check(session, req)
	}

	return nil
}

//AutoGet will get session, if no vaild session
//it will create a new session,and reurn second value as false
func AutoGet(sessionID string, req *http.Request) (SessionIO, bool) {

	if session, ok := sessionMap[sessionID]; ok {
		if session := check(session, req); session != nil {
			return session, true
		}
		return New(req), false
	}

	return New(req), false

}

//check wether session experid and wether belong to req
func check(session SessionIO, req *http.Request) SessionIO {
	if !session.expired() && session.belong(req) {
		return session
	}
	return nil
}

func checkExpire(sessionMap map[string]SessionIO) {
	for true {
		for sessionID, session := range sessionMap {
			if session.expired() {
				//此处无需加锁，session被移除不会影响已经取出的session的工作。
				delete(sessionMap, sessionID)
			}
			time.Sleep(time.Second)
		}
	}
}
