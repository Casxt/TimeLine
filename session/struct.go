package session

//Session struct
type Session struct {
	sessionID  string
	expireTime int
	Map        map[string]interface{}
}

func (session Session) Get(key string) string {
	if res, ok := session.Map[key].(string); ok {
		return res
	} else {
		return ""
	}
}
