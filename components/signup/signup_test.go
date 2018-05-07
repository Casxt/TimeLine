package signup

import (
	"testing"

	"github.com/Casxt/TimeLine/database"
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
