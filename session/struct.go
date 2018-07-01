package session

import (
	"net/http"
	"sync"
	"time"
)

//SessionObj struct should not be changed, readonly
type sessionObj struct {
	value      interface{}
	expireTime time.Duration
	setTime    time.Time
}

func (obj *sessionObj) expired() bool {
	if time.Since(obj.setTime) > obj.expireTime && obj.expireTime != 0 {
		return false
	}
	return true
}

//Session struct
type Session struct {
	sessionID  string
	expireTime time.Duration
	address    string
	userAgent  string
	setTime    time.Time
	Map        map[string]*sessionObj
	lock       sync.RWMutex
}

func (session *Session) init(sessionID string) {
	session.Lock()
	defer session.Unlock()

	session.sessionID = sessionID
	session.Map = make(map[string]*sessionObj)
}

//ExtraInfo set ExtraInfo for auth,
//set "" will not change, can use to get ExtraInfo
func (session *Session) ExtraInfo(address, userAgent string) (string, string) {
	if address != "" {
		session.Lock()
		session.address = address
		session.Unlock()
	}
	if userAgent != "" {
		session.Lock()
		session.userAgent = userAgent
		session.Unlock()
	}
	session.RLock()
	defer session.RUnlock()
	return session.address, session.userAgent
}

//ExpireTime set the ExpireTime of session and retuen ExpireTime of session,
//ExpireTime smaller than 0 will not change ExpireTime of session so than can be use to get ExpireTime
func (session *Session) ExpireTime(expireTime time.Duration) time.Duration {
	session.refresh()

	session.RLock()
	defer session.RUnlock()

	if expireTime < 0 {
		return session.expireTime
	}
	session.expireTime = expireTime
	return session.expireTime
}

//Belong check wether req match the requirement of session
func (session *Session) Belong(req *http.Request) bool {
	if req.RemoteAddr == session.address ||
		req.UserAgent() == session.userAgent {
		return true
	}
	return false
}

//ID return sessionID
func (session *Session) ID() string {
	return session.sessionID
}

//Expired return wether session is expired
func (session *Session) Expired() bool {
	session.RLock()
	defer session.RUnlock()

	if time.Since(session.setTime) > session.expireTime {
		session.Map = nil
		return true
	}
	return false
}

func (session *Session) refresh() {
	session.Lock()
	session.setTime = time.Now()
	session.Unlock()
}

//Get  return type is string
//if not have key or value type is not string return ""
func (session *Session) Get(key string) (res string, ok bool) {
	session.refresh()
	session.RLock()
	if session.Map == nil {
		return "", false
	}
	Obj, ok := session.Map[key]
	if !ok {
		return "", false
	}
	session.RUnlock()
	if Obj.expired() {
		session.Delete(key)
		return "", false
	}
	if res, ok := Obj.value.(string); ok {
		return res, true
	}
	return "", false
}

//GetInt return type is int
func (session *Session) GetInt(key string) (res int, ok bool) {
	session.refresh()
	session.RLock()
	if session.Map == nil {
		return 0, false
	}
	Obj, ok := session.Map[key]
	if !ok {
		return 0, false
	}
	session.RUnlock()
	if Obj.expired() {
		session.Delete(key)
		return 0, false
	}
	res, ok = Obj.value.(int)
	return res, ok
}

//GetTime return type is time.Time
func (session *Session) GetTime(key string) (res time.Time, ok bool) {
	session.refresh()
	session.lock.RLock()
	if session.Map == nil {
		return time.Time{}, false
	}
	Obj, ok := session.Map[key]
	if !ok {
		return time.Time{}, false
	}
	session.lock.RUnlock()
	if Obj.expired() {
		session.Delete(key)
		return time.Time{}, false
	}
	res, ok = Obj.value.(time.Time)
	return res, ok
}

//Have if session have key return true
func (session *Session) Have(key string) bool {
	session.refresh()
	session.lock.RLock()
	_, ok := session.Map[key]
	session.lock.RUnlock()
	return ok
}

//Put string
//set expireTime=0 to enable infinit lifetime
func (session *Session) Put(key string, value string, expireTime time.Duration) {
	session.refresh()
	session.lock.Lock()
	defer session.lock.Unlock()
	session.Map[key] = &sessionObj{value: value, expireTime: expireTime, setTime: time.Now()}
}

//PutInt int
//set expireTime=0 to enable infinit lifetime
func (session *Session) PutInt(key string, value int, expireTime time.Duration) {
	session.refresh()
	session.lock.Lock()
	defer session.lock.Unlock()
	session.Map[key] = &sessionObj{value: value, expireTime: expireTime, setTime: time.Now()}
}

//PutTime time.Time
//set expireTime=0 to enable infinit lifetime
func (session *Session) PutTime(key string, value time.Time, expireTime time.Duration) {
	session.refresh()
	session.lock.Lock()
	defer session.lock.Unlock()
	session.Map[key] = &sessionObj{value: value, expireTime: expireTime, setTime: time.Now()}
}

//Delete key,
//delete will not cause panic
func (session *Session) Delete(key string) {
	session.Lock()
	delete(session.Map, key)
	session.Unlock()
}

//RLock RLock
func (session *Session) RLock() {
	session.lock.RLock()
}

//RUnlock RUnlock
func (session *Session) RUnlock() {
	session.lock.RUnlock()
}

//Lock Lock
func (session *Session) Lock() {
	session.lock.Lock()
}

//Unlock Unlock
func (session *Session) Unlock() {
	session.lock.Unlock()
}
