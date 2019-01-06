package tools

import (
	"encoding/json"
	"github.com/Casxt/TimeLine/config"
	"net/http"
	"strings"
)

//GetWeiXinUser
func GetWeiXinUser(code string) (OpenID, SessionKey string, UnionID *string, status int, jsonRes map[string]string) {
	resp, err := http.Get(strings.Join([]string{`https://api.weixin.qq.com/sns/jscode2session?appid=`, config.WeiXinApp.Id,
		`&secret=`, config.WeiXinApp.Secrete,
		`&js_code=`, code, `&grant_type=authorization_code`}, ""))
	if err != nil {
		return "", "", nil, 400, map[string]string{
			"State": "Failed",
			"Msg":   `Connect to weixin Server Failed`,
		}
	}
	if resp.StatusCode != 200 {

	}
	type UserData struct {
		OpenId     string  `json:"openid"`
		SessionKey string  `json:"session_key"`
		UnionId    *string `json:"unionid,omitempty"`
		ErrorCode  int     `json:"errcode"`
		ErrorMsg   int     `json:"errmsg"`
	}
	var userData UserData

	if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
		return "", "", nil, 400, map[string]string{
			"State": "Failed",
			"Msg":   "Invalid Json when auth with WeiXin",
		}
	}
	switch userData.ErrorCode {
	case 0:
		return userData.OpenId, userData.SessionKey, userData.UnionId, 200, map[string]string{
			"State": "Failed",
			"Msg":   "Invalid Json when auth with WeiXin",
		}
	case -1:
		return userData.OpenId, userData.SessionKey, userData.UnionId, 500, map[string]string{
			"State": "Failed",
			"Msg":   "WeiXin server are too busy, please try again",
		}
	case 40029:
		return userData.OpenId, userData.SessionKey, userData.UnionId, 400, map[string]string{
			"State": "Failed",
			"Msg":   "Invalid UserCode",
		}
	case 45011:
		return userData.OpenId, userData.SessionKey, userData.UnionId, 400, map[string]string{
			"State": "Failed",
			"Msg":   "request WeiXin server are too fast, please try later",
		}
	default:
		return userData.OpenId, userData.SessionKey, userData.UnionId, 400, map[string]string{
			"State": "Failed",
			"Msg":   "WeiXin server return unknown code",
		}
	}
}
