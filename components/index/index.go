package index

import (
	"net/http"

	"github.com/Casxt/TimeLine/static"
)

//Route decide static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	status, result, _ = static.GetPage("components", "index", "index.html")
	res.WriteHeader(status)
	res.Write(result)
}
