package profile

import (
	"net/http"

	"github.com/Casxt/TimeLine/static"
)

//Route decide static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/profile"):]
	switch {
	default:
		status, result, _ = static.GetPage("components", "profile", "profile.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}
