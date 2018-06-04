package session

import (
	"net/http"
	"time"
)

//IO Interface
type IO interface {
	Manager
	Reader
	Writer
}

//Manager Interface
type Manager interface {
	init(sessionID string)
	ExpireTime(expireTime time.Duration) time.Duration
	Belong(*http.Request) bool
	RLock()
	RUnlock()
	Lock()
	Unlock()
	Expired() bool
	ID() string
	refresh()
}

//Writer Interface
type Writer interface {
	Put(key, value string)
	PutInt(key string, value int)
	PutTime(key string, value time.Time)

	PutAll(sessionKV map[string]interface{})
}

//Reader Interface
type Reader interface {
	Get(key string) (string, bool)
	GetInt(key string) (int, bool)
	GetTime(key string) (time.Time, bool)

	GetAll() map[string]interface{}
}
