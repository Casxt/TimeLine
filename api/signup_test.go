package api

import (
	"testing"

	"github.com/MapleFadeAway/timeline/database"
)

func TestSignUp(t *testing.T) {
	database.Open()
	defer database.Close()
	/*res := SignUp([]byte(`{
	"Phone":"18110020002",
	"Mail":"774714620@qq.com",
	"HashPass":""}`))*/
	//log.Println(res)
}
