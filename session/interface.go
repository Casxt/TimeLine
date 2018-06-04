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
	ExtraInfo(address, userAgent string) (string, string)
	RLock()
	RUnlock()
	Lock()
	Unlock()
	Expired() bool
	ID() string
	refresh()
	Delete()
}

//Writer Interface
type Writer interface {
	Put(key, value string, expireTime time.Duration)
	PutInt(key string, value int, expireTime time.Duration)
	PutTime(key string, value time.Time, expireTime time.Duration)
}

//Reader Interface
type Reader interface {
	Get(key string) (string, bool)
	GetInt(key string) (int, bool)
	GetTime(key string) (time.Time, bool)
}
