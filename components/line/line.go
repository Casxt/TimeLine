package line

import (
	"net/http"
	"strings"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/page"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide page
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	subPath := req.URL.Path[len("/line"):]
	switch {
	case strings.HasPrefix(strings.ToLower(subPath), "create"):
		result, status, _ = page.GetPage("components", "line", "createLine.html")
	default:
		result, status, _ = page.GetPage("components", "line", "line.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

//CreateLine will create a new line with specific name
//TODO: Limit User num of Line
func CreateLine(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		Operator string `json:"Operator"`
		LineName string `json:"LineName"`
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, session := tools.GetLoginState(req, data.Operator)
	if session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   `User Not Sign In`,
		}
	}

	err := database.CreateLine(data.LineName, UserID)
	if err != nil {
		switch err.Error() {
		case "Line Already Exist":
			return 400, map[string]string{
				"State": "Failde",
				"Msg":   "Line Name Already be Used",
			}

		default:
			return 400, map[string]string{
				"State": "Failde",
				"Msg":   "Line Name Already be Used",
			}
		}
	}
	return 200, map[string]string{
		"State": "Success",
		"Msg":   "Line create success",
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
