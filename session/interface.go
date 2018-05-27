package session

import (
	"time"
)

//SessionIO Interface
type SessionIO interface {
	Manager
	Reader
	Writer
}

//Manager Interface
type Manager interface {
	Init(sessionID string)
	ExpireTime(expireTime time.Duration) time.Duration
	RLock()
	RUnlock()
	Lock()
	Unlock()
	expired() bool
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
