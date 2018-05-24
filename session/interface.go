package session

import (
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
	New(sessionID string) *Session
	ExpireTime(time int) int
}

//Writer Interface
type Writer interface {
	Put(key, value string) error
	PutInt(key string, value int) error
	PutTime(key string, value time.Time) error

	PutAll(sessionKV map[string]string) error
}

//Reader Interface
type Reader interface {
	Get(key string) string
	GetInt(key string) int
	GetTime(key string) time.Time

	GetAll() map[string]string
}
