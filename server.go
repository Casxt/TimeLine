package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/Casxt/TimeLine/api"
	"github.com/Casxt/TimeLine/components/image"
	"github.com/Casxt/TimeLine/components/line"
	"github.com/Casxt/TimeLine/components/profile"
	"github.com/Casxt/TimeLine/components/signin"
	"github.com/Casxt/TimeLine/components/signup"
	"github.com/Casxt/TimeLine/database"
	"github.com/Casxt/TimeLine/static"
)

func route(res http.ResponseWriter, req *http.Request) {
	path := req.URL.Path

	log.Println(req.RemoteAddr[0:strings.LastIndex(req.RemoteAddr, ":")], req.Method, path)
	switch {
	case strings.HasPrefix(strings.ToLower(path), "/api"):
		api.Route(res, req)
	case strings.HasPrefix(strings.ToLower(path), "/signin"):
		signin.Route(res, req)
	case strings.HasPrefix(strings.ToLower(path), "/signup"):
		signup.Route(res, req)
	case strings.HasPrefix(strings.ToLower(path), "/static"):
		static.Route(res, req)
	case strings.HasPrefix(strings.ToLower(path), "/image"):
		image.Route(res, req)
	case strings.HasPrefix(strings.ToLower(path), "/line"):
		line.Route(res, req)
	case strings.HasPrefix(strings.ToLower(path), "/profile"):
		profile.Route(res, req)
	default:
		//line.Route(res, req)
	}
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
