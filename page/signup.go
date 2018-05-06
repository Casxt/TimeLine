package page

import "net/http"

//Route deliver the req depend on method
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int
	switch req.Method {
	case "GET":
		status, result = getSignUpPage()
	default:
		status, result = getSignUpPage()
	}
	res.WriteHeader(status)
	res.Write(result)
}

func getSignUpPage() (status int, res []byte) {
	res, err := Get("signup.html")
	status = 200
	if err != nil {
		status = 404
		res, _ = Get("404.html")
	}
	return status, res
}
