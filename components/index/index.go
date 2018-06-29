package index

import (
	"net/http"

	"github.com/Casxt/TimeLine/page"
)

func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/index"):]
	switch {
	//case strings.HasPrefix(strings.ToLower(subPath), "signin.js"):
	//	result, status, _ = page.GetFile("components", "signin", "signup.js")
	default:
		result, status, _ = page.GetPage("components", "index", "index.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}
