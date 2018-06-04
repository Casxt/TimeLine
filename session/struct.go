package session

import (
	"net/http"
	"sync"
	"time"
)

//Session struct
type Session struct {
	sessionID  string
	expireTime time.Duration
	address    string
	setTime    time.Time
	Map        map[string]interface{}
	lock       sync.RWMutex
}

func (session *Session) init(sessionID string) {
	session.lock.Lock()
	defer session.lock.Unlock()

	session.sessionID = sessionID
	session.Map = make(map[string]interface{})

}

//ExpireTime set the ExpireTime of session and retuen ExpireTime of session,
//ExpireTime smaller than 0 will not change ExpireTime of session so than can be use to get ExpireTime
func (session *Session) ExpireTime(expireTime time.Duration) time.Duration {
	session.lock.RLock()
	defer session.lock.RUnlock()

	if expireTime < 0 {
		return session.expireTime
	}

	session.lock.Lock()
	defer session.lock.Unlock()
	session.setTime = time.Now()
	session.expireTime = expireTime
	return session.expireTime
}

func (session *Session) Belong(req *http.Request) bool {
	if req.RemoteAddr == session.address {
		return true
	}
	return false
}

func (session *Session) ID() string {
	return session.sessionID
}

func (session *Session) Expired() bool {
	session.lock.RLock()
	defer session.lock.RUnlock()

	if time.Since(session.setTime) > session.expireTime {
		session.Map = nil
		return false
	}
	return true
}

func (session *Session) refresh() {
	session.lock.Lock()
	defer session.lock.Unlock()
	session.setTime = time.Now()
}

func (session *Session) Get(key string) (res string, ok bool) {
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session.Map == nil {
		return "", false
	}
	session.refresh()
	res, ok = session.Map[key].(string)
	return res, ok
}

func (session *Session) GetInt(key string) (res int, ok bool) {
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session.Map == nil {
		return 0, false
	}
	session.refresh()
	res, ok = session.Map[key].(int)
	return res, ok
}

func (session *Session) GetTime(key string) (res time.Time, ok bool) {
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session.Map == nil {
		return time.Time{}, false
	}
	session.refresh()
	res, ok = session.Map[key].(time.Time)
	return res, ok
}

func (session *Session) GetAll() map[string]interface{} {
	session.lock.RLock()
	defer session.lock.RUnlock()
	if session.Map == nil {
		return nil
	}
	session.refresh()
	return session.Map
}

func (session *Session) Put(key string, value string) {
	session.lock.Lock()
	defer session.lock.Unlock()
	if session.Map == nil {
		session.Map = make(map[string]interface{})
	}
	session.refresh()
	session.Map[key] = value
}

func (session *Session) PutInt(key string, value int) {
	session.lock.Lock()
	defer session.lock.Unlock()
	if session.Map == nil {
		session.Map = make(map[string]interface{})
	}
	session.refresh()
	session.Map[key] = value
}

func (session *Session) PutAll(Map map[string]interface{}) {
	session.lock.Lock()
	defer session.lock.Unlock()
	session.refresh()
	session.Map = Map
}

func (session *Session) PutTime(key string, value time.Time) {
	session.lock.Lock()
	defer session.lock.Unlock()
	if session.Map == nil {
		session.Map = make(map[string]interface{})
	}
	session.refresh()
	session.Map[key] = value
}

func (session *Session) RLock() {
	session.lock.RLock()
}

func (session *Session) RUnlock() {
	session.lock.RUnlock()
}

func (session *Session) Lock() {
	session.lock.Lock()
}

func (session *Session) Unlock() {
	session.lock.Unlock()
}
