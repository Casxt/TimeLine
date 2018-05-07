package page

import (
	"net/http"
	"strings"
)

//Route decide return which Page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	//page路由
	subPath := req.URL.Path[8:] //url are like /static/(something)
	switch {
	case strings.HasPrefix(strings.ToLower(subPath), "js"):
		result, status, _ = GetFile(strings.Split(subPath, "/")...)
	}
	res.WriteHeader(status)
	res.Write(result)
}
