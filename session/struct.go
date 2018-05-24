package session

import (
	"time"
)

//Session struct
type Session struct {
	sessionID  string
	expireTime time.Duration
	setTime    time.Time
	Map        map[string]interface{}
}

func (session Session) ExpireTime(expireTime time.Duration) time.Duration {
	if expireTime == 0 {
		return session.expireTime
	}

	session.setTime = time.Now()
	session.expireTime = expireTime
	return session.expireTime
}

func (session Session) expired() bool {
	if time.Since(session.setTime) > session.expireTime {
		session.Map = nil
		return false
	}
	return true
}

func (session Session) refresh() {
	session.setTime = time.Now()
}

func (session Session) Get(key string) (res string, ok bool) {
	if session.expired() {
		return "", false
	}
	session.refresh()
	res, ok = session.Map[key].(string)
	return res, ok
}

func (session Session) GetInt(key string) (res int, ok bool) {
	if session.expired() {
		return 0, false
	}
	session.refresh()
	res, ok = session.Map[key].(int)
	return res, ok
}

func (session Session) GetTime(key string) (res time.Time, ok bool) {
	if session.expired() {
		return time.Time{}, false
	}
	session.refresh()
	res, ok = session.Map[key].(time.Time)
	return res, ok
}

func (session Session) GetAll() map[string]interface{} {
	if session.expired() {
		return nil
	}
	session.refresh()
	return session.Map
}

func (session Session) Put(key string, value string) {
	if session.expired() {
		session.Map = make(map[string]interface{})
	}
	session.refresh()
	session.Map[key] = value
}

func (session Session) PutInt(key string, value int) {
	if session.expired() {
		session.Map = make(map[string]interface{})
	}
	session.refresh()
	session.Map[key] = value
}

func (session Session) PutAll(Map map[string]interface{}) {
	//if session.expired() {
	//session.Map = make(map[string]interface{})
	//}
	session.refresh()
	session.Map = Map
}

func (session Session) PutTime(key string, value time.Time) {
	if session.expired() {
		session.Map = make(map[string]interface{})
	}
	session.refresh()
	session.Map[key] = value
}
