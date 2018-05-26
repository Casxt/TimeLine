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
	session := new(Session)
	session.ExpireTime(time.Duration(time.Hour * 24 * 30))
	session.Map["RemoteAddr"] = req.RemoteAddr
	session.sessionID = newID()
	count := 0

	for sessionMap[session.sessionID] != nil && count < 3 {
		session.sessionID = newID()
		count++
	}

	if count == 3 {
		panic("too many key duplicate")
	}

	sessionMap[session.sessionID] = session
	return session
}

func Get(sessionID string, req *http.Request) SessionIO {
	return check(sessionMap[sessionID], req)

}

func check(session SessionIO, req *http.Request) SessionIO {
	if session.expired() {
		return nil
	}
	if addr, _ := session.Get("RemoteAddr"); addr == req.RemoteAddr {
		return nil
	}
	return session
}
