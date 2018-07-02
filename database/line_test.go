package database

import "testing"

func TestCreateLine(t *testing.T) {
	Open()
	err := CreateLine("testline", "116")
	if err != nil {
		t.Error("CreateLine Error", err.Error())
	}
	t.Log("success")
}
