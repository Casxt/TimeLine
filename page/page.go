package page

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
)

//ProjectPath is path that save the static file
var ProjectPath = "C:\\Users\\Surface\\go\\src\\github.com\\Casxt\\TimeLine"

//GetPage return the builted html
func GetPage(names ...string) (Page []byte, err error) {
	var absPath string
	var content []byte

	absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "Header.html"))
	content, err = ioutil.ReadFile(absPath)
	buf := bytes.NewBuffer(content)
	if err != nil {
		return nil, err
	}

	absPath = filepath.FromSlash(path.Join(ProjectPath, path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		return nil, err
	}

	absPath = filepath.FromSlash(path.Join(ProjectPath, "static", "Footer.html"))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}

//GetFile Get static file
func GetFile(names ...string) (File []byte, status int, err error) {
	var absPath string
	var content []byte

	absPath = filepath.FromSlash(path.Join(ProjectPath, "static", path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)

	if err != nil {
		return nil, 500, err
	}

	return content, 200, err
}
