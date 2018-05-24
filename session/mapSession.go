package session

var (
	sessionMap map[string]*Session
)

func Open() {
	sessionMap = make(map[string]*Session)
}
