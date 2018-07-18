package session

import (
	"net/http"
	"strings"
	"time"
)

var (
	sessionMap map[string]IO
)

//Open session server
func Open() {
	sessionMap = make(map[string]IO)
	go checkExpire(sessionMap)
}

//New session
func New(req *http.Request) IO {
	var session IO
	var sessionID string
	sessionID = newID()
	count := 0

	// in case of key duplicate but will it really happend?
	for sessionMap[sessionID] != nil && count < 3 {
		sessionID = newID()
		count++
	}

	if count == 3 {
		panic("too many key duplicate")
	}

	session = new(Session)
	session.init(sessionID)
	session.ExpireTime(time.Duration(time.Hour * 24 * 30))
	session.ExtraInfo(req.RemoteAddr[0:strings.LastIndex(req.RemoteAddr, ":")], req.UserAgent())
	sessionMap[sessionID] = session
	return session
}

//GetByCookie get session by cookie
func GetByCookie(req *http.Request) IO {
	var cookie *http.Cookie
	var err error
	if cookie, err = req.Cookie("SessionID"); err != nil {
		return nil
	}
	return Get(cookie.Value, req)
}

//Get and check session
func Get(sessionID string, req *http.Request) IO {
	if session, ok := sessionMap[sessionID]; ok {
		if check(session, req) {
			return session
		}
		return nil
	}
	return nil
}

//Auto will check-get session or create new one,
//try to find sessionID in req cookie
//if no vaild session
//it will create a new session,and reurn second value as true
//in addition this func will add cookie
func Auto(res http.ResponseWriter, req *http.Request) (session IO, NewSession bool) {

	var cookie *http.Cookie
	var err error
	if cookie, err = req.Cookie("SessionID"); err != nil {
		session = New(req)
		http.SetCookie(res, &http.Cookie{Name: "SessionID", Value: session.ID(), Path: "/", MaxAge: 30 * 86400}) //, HttpOnly: true
		return session, true
	}
	sessionID := cookie.Value

	if session, ok := sessionMap[sessionID]; ok {
		if check(session, req) {
			http.SetCookie(res, &http.Cookie{Name: "SessionID", Value: session.ID(), Path: "/", MaxAge: 30 * 86400}) //, HttpOnly: true
			return session, false
		}
		// session extra info not match, create a new one
		session = New(req)
		http.SetCookie(res, &http.Cookie{Name: "SessionID", Value: session.ID(), Path: "/", MaxAge: 30 * 86400}) //, HttpOnly: true
		return session, true
	}
	session = New(req)
	http.SetCookie(res, &http.Cookie{Name: "SessionID", Value: session.ID(), Path: "/", MaxAge: 30 * 86400}) //, HttpOnly: true
	return session, true

}

//check wether session Expired and wether belong to req
func check(session IO, req *http.Request) bool {
	if !session.Expired() && session.Belong(req) {
		return true
	}
	//if delete here, may cause error
	return false
}

//checkExpire will delete expeired session
//need to optimal range sessionMap
func checkExpire(sessionMap map[string]IO) {
	var counter int8
	for true {
		for sessionID, session := range sessionMap {
			counter++
			if session.Expired() {
				//log.Panicln("Delete expeired Session:", sessionID)
				delete(sessionMap, sessionID)
				//考虑主动释放
			}
			if counter > 9 {
				time.Sleep(time.Second)
				counter = 0
			}
		}
		//if sessionMap is empty, without this will cause not paused forever loop
		time.Sleep(time.Second)
	}
}
