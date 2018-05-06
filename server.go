package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/MapleFadeAway/timeline/api"
	"github.com/MapleFadeAway/timeline/database"
	"github.com/MapleFadeAway/timeline/page"
)

func route(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	switch {
	case strings.HasPrefix(path, "/api"):
		api.Route(res, req)
	case strings.HasPrefix(path, "/SignUp"):
		page.Route(res, req)
	default:
		//不要为首页创建专门的判断，所有的首页都应该被默认展示
		res.Write([]byte("TimeLine!"))
	}
}

func run() {

	if err := database.Open(); err != nil {
		log.Println("database Open filed")
		return
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(route))

	log.Println("Listening...")
	http.ListenAndServe(":80", mux)
}

//SetUp will setup database table
func SetUp() {
	if err := database.Open(); err != nil {
		log.Println("database Open filed")
		return
	}
	defer database.Close()

	course, selfCourse, err := database.Begin(nil)
	if err != nil {
		log.Fatalln(err)
		return
	}

	if err := database.CreateUserTable(course); err != nil {
		database.Rollback(course, selfCourse)
		return
	}
	if err := database.CreateLineTable(course); err != nil {
		database.Rollback(course, selfCourse)
		return
	}
	if err := database.CreateGroupTable(course); err != nil {
		database.Rollback(course, selfCourse)
		return
	}
	if err := database.CreateSliceTable(course); err != nil {
		database.Rollback(course, selfCourse)
		return
	}

	database.Commit(course, selfCourse)
}
