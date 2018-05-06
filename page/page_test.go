package page

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	res, err := Get(`index.html`)
	log.Fatalln(string(res))
	if err != nil {
		t.Error("TestGet Error", err)
	}
}
