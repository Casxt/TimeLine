package session

import (
	"net/http"
	"sync"
	"time"
)

var (
	sessionMap map[string]SessionIO
)

func Open() {
	sessionMap = make(map[string]SessionIO)
	mapLock := new(sync.RWMutex)
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

	session.Lock()
	session.Init(sessionID)
	session.ExpireTime(time.Duration(time.Hour * 24 * 30))
	session.Put("RemoteAddr", req.RemoteAddr)
	session.Unlock()

	sessionMap[sessionID] = session
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

func checkExpire(sessionMap map[string]SessionIO) {
	for true {
		for sessionID, session := range sessionMap {
			if session.expired() {
				delete(sessionMap, sessionID)
				time.Sleep(time.Second)
			}
		}
	}
}
