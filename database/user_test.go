package database

import (
	"log"
	"testing"
)

//CreateUser create a unverify user
func TestCreateUser(t *testing.T) {
	Open()
	NickName, Pass, err := CreateUser("18000000000", "770770770", "")
	if err != nil {
		t.Error("TestCreateUser Error", err)
	}
	log.Println(NickName, Pass)
}

//TestGetUserByPhone
func TestGetUserByPhone(t *testing.T) {
	Open()
	ID, Mail, Pass, Gender, Salt, SaltPass, ProfilePic, SignInTime, err := GetUserByPhone("18110020001")
	if err != nil {
		t.Error("TestGetUserByPhone Error", err)
	}
	t.Log(ID, Mail, Pass, Gender, Salt, SaltPass, ProfilePic, SignInTime)
}
