package session

import (
	"time"
)

var (
	sessionMap map[string]SessionIO
)

func Open() {
	sessionMap = make(map[string]SessionIO)
}

func New() SessionIO {
	session := new(Session)
	session.ExpireTime(time.Duration(time.Hour * 24 * 30))
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
