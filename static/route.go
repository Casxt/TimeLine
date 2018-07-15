package static

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
	case strings.HasPrefix(strings.ToLower(subPath), "js"): // /(js/...)
		status, result, _ = GetFile(strings.Split(req.URL.Path, "/")...)
	case strings.HasPrefix(strings.ToLower(subPath), "css"): // /(js/...)
		status, result, _ = GetFile(strings.Split(req.URL.Path, "/")...)
	}
	res.WriteHeader(status)
	res.Write(result)
}
