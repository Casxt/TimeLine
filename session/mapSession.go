package session

import (
	"net/http"
	"time"
)

var (
	sessionMap map[string]IO
)

func Open() {
	sessionMap = make(map[string]IO)
	go checkExpire(sessionMap)
}

func New(req *http.Request) IO {
	var session IO
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

	session.init(sessionID)
	session.ExpireTime(time.Duration(time.Hour * 24 * 30))
	session.Put("RemoteAddr", req.RemoteAddr)

	sessionMap[sessionID] = session
	return session
}

func Get(sessionID string, req *http.Request) IO {

	if session, ok := sessionMap[sessionID]; ok {
		return check(session, req)
	}

	return nil
}

//Auto will get session,
//if no vaild session
//it will create a new session,and reurn second value as false
func Auto(sessionID string, req *http.Request) (IO, bool) {

	if session, ok := sessionMap[sessionID]; ok {
		if session := check(session, req); session != nil {
			return session, true
		}
		// session extra info not match, create a new one
		session := New(req)
		req.AddCookie(&http.Cookie{Name: "SessionId", Value: session.ID(), Path: "/", MaxAge: 86400})
		return session, false
	}

	return New(req), false

}

//check wether session Expired and wether belong to req
func check(session IO, req *http.Request) IO {
	if !session.Expired() && session.Belong(req) {
		return session
	}
	return nil
}

func checkExpire(sessionMap map[string]IO) {
	for true {
		for sessionID, session := range sessionMap {
			if session.Expired() {
				//此处无需加锁，session被移除不会影响已经取出的session的工作。
				delete(sessionMap, sessionID)
			}
			time.Sleep(time.Second)
		}
	}
}
