package signin

import "net/http"

//Route Return The Page to Show
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	switch req.Method {
	case "GET":
		//status, result = Page()
	default:
		//status, result = Page()
	}
	res.WriteHeader(status)
	res.Write(result)
}

//
func CheckAccount() {

}
