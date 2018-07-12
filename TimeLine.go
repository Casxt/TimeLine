package main

import (
	"log"
	"net/http"

	"github.com/Casxt/TimeLine/config"
	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/session"
)

func main() {
	config.Load("config.json")

	log.Println("Connect Database...")
	if err := database.Open(); err != nil {
		log.Println("database Open filed")
		return
	}
	defer database.Close()

	log.Println("Start Session server...")
	session.Open()
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(route))

	log.Println("Listening...")
	http.ListenAndServe(":80", mux)
}
