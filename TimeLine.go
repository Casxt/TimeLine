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

	if config.TLS.Cert != "" && config.TLS.Key != "" {
		log.Println("TLS Enable")
		log.Println("Start Https Server @ 443 ...")
		if err := http.ListenAndServeTLS(":443", config.TLS.Cert, config.TLS.Key, mux); err != nil {
			log.Fatalln(err.Error())
		}
	} else {
		log.Println("TLS Disable")
		log.Println("Start Http Server @ 80 ...")
		if err := http.ListenAndServe(":80", mux); err != nil {
			log.Fatalln(err.Error())
		}
	}
}
