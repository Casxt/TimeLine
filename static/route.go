package static

import (
	"mime"
	"net/http"
	"path"
	"strings"
)

//Route decide return which Page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	//page路由
	subPath := req.URL.Path[len("/static"):] //url are like /static/(something)
	splitPath := strings.Split(req.URL.Path, "/")
	switch {
	case strings.HasPrefix(strings.ToLower(subPath), "/js"):
		status, result, _ = GetFile(splitPath...)
	case strings.HasPrefix(strings.ToLower(subPath), "/css"):
		status, result, _ = GetFile(splitPath...)
	}
	res.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(splitPath[len(splitPath)-1])))
	res.WriteHeader(status)
	res.Write(result)
}
