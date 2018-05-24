package page

import (
	"bytes"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
)

//ProjectPath is path that save the static file
var ProjectPath = "C:\\Users\\Surface\\go\\src\\github.com\\Casxt\\TimeLine"

//GetPage return the builted html
func GetPage(names ...string) (Page []byte, status int, err error) {
	var absPath string
	var content []byte
	var buf *bytes.Buffer

	absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "Header.html"))
	content, err = ioutil.ReadFile(absPath)
	buf = bytes.NewBuffer(content)
	if err != nil {
		absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "404.html"))
		content, err = ioutil.ReadFile(absPath)
		buf = bytes.NewBuffer(content)
		return buf.Bytes(), 404, err
	}

	absPath = filepath.FromSlash(path.Join(ProjectPath, path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "404.html"))
		content, err = ioutil.ReadFile(absPath)
		buf = bytes.NewBuffer(content)
		return buf.Bytes(), 404, err
	}

	absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "Footer.html"))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "404.html"))
		content, err = ioutil.ReadFile(absPath)
		buf = bytes.NewBuffer(content)
		return buf.Bytes(), 404, err
	}

	return buf.Bytes(), 200, err
}

//GetFile Get static file
func GetFile(names ...string) (File []byte, status int, err error) {
	var absPath string
	var content []byte
	//res http.ResponseWriter,
	absPath = filepath.FromSlash(path.Join(ProjectPath, path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)

	if err != nil {
		log.Println(err.Error())
		return nil, 500, err
	}
	//log.Println(mime.TypeByExtension(path.Ext(names[len(names)-1])))
	//log.Println(path.Ext(names[len(names)-1]))
	//res.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(names[len(names)-1])))
	return content, 200, err
}
