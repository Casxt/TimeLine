package line

import (
	"net/http"

	"github.com/Casxt/TimeLine/page"
)

//Route decide page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	//subPath := req.URL.Path[len("/index"):]
	switch {
	//case strings.HasPrefix(strings.ToLower(subPath), "signin.js"):
	//	result, status, _ = page.GetFile("components", "signin", "signup.js")
	default:
		result, status, _ = page.GetPage("components", "line", "line.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

//CreateLine will create a new line with specific name
func CreateLine(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		Operator string `json:"Operator"`
		Name     string `json:"Name"`
	}

	return 400, map[string]string{
		"State": "Failde",
		"Msg":   "Invilde Parameter",
	}
}

//AddUser will will add user to specific line
func AddUser(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		Operator  string `json:"Operator"`
		Name      string `json:"Name"`
		UserPhone string `json:"UserPhone"`
	}

	return 400, map[string]string{
		"State": "Failde",
		"Msg":   "Invilde Parameter",
	}
}
