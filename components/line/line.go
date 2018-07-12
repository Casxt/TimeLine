package line

import (
	"net/http"
	"strings"

	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/static"
	"github.com/Casxt/TimeLine/tools"
)

//Route decide static
func Route(res http.ResponseWriter, req *http.Request) {
	var result []byte
	var status int

	subPath := req.URL.Path[len("/line"):]
	switch {
	//match /linecreate*
	case strings.HasPrefix(strings.ToLower(subPath), "create"):
		status, result, _ = static.GetPage("components", "line", "createLine.html")
	default:
		status, result, _ = static.GetPage("components", "line", "line.html")
	}
	res.WriteHeader(status)
	res.Write(result)
}

//GetLines Get Lines of User
func GetLines(res http.ResponseWriter, req *http.Request) (status int, jsonRes interface{}) {
	type Data struct {
		Operator  string
		SessionID string
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, _ := tools.GetLoginStateOfOperator(req, data.SessionID, data.Operator)
	if UserID == "" {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   `User Not Sign In`,
		}
	}

	Lines, err := database.GetLines(UserID)
	if err != nil {
		return 500, map[string]string{
			"State":  "Failde",
			"Msg":    `database.GetUserLines Failed`,
			"Detial": err.Error(),
		}
	}

	type ResData struct {
		State string
		Msg   string
		Lines []string
	}

	return 200, ResData{
		State: "Success",
		Msg:   "Lines Get Successful",
		Lines: Lines,
	}
}

//CreateLine will create a new line with specific name
//TODO: Limit User num of Line
func CreateLine(res http.ResponseWriter, req *http.Request) (status int, jsonRes map[string]string) {
	type Data struct {
		Operator  string
		LineName  string
		SessionID string
	}
	var data Data

	if status, jsonRes = tools.GetPostJSON(req, &data); status != 200 {
		return status, jsonRes
	}

	UserID, session := tools.GetLoginStateOfOperator(req, data.SessionID, data.Operator)
	if session == nil {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   `User Not Sign In`,
		}
	}

	if len(data.LineName) < 4 {
		return 400, map[string]string{
			"State": "Failde",
			"Msg":   "Line Name too short",
		}
	}

	err := database.CreateLine(strings.ToLower(data.LineName), UserID)
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
