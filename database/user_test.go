package database

import (
	"testing"
)

//CreateUser create a unverify user
func TestCreateUser(t *testing.T) {
	Open()
	if err := CreateUser("18000000000", "770770770"); err != nil {
		t.Error("TestCreateUser Error", err)
	}
}
