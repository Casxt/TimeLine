package session

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

//Md5 caculate the md5 of string
//use to generate id
func Md5(text string) string {
	hashMd5 := md5.New()
	hashMd5.Write([]byte(text))
	return hex.EncodeToString(hashMd5.Sum(nil))
}

//newID return a new session Id
func newID() (sessionID string) {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	sessionID = Md5(strconv.FormatInt(nano, 10)) + Md5(strconv.FormatInt(rndNum, 10))
	return sessionID
}
