package page

import (
	"bytes"
	"io/ioutil"
	"path"
	"path/filepath"
)

//StaticPath is path that save the static file
var StaticPath = "C:\\Users\\Surface\\go\\src\\github.com\\Casxt\\TimeLine\\static"

//GetPage return the builted html
func GetPage(name string) (Page []byte, err error) {
	var absPath string
	var content []byte

	absPath = filepath.FromSlash(path.Join(StaticPath, "Header.html"))
	content, err = ioutil.ReadFile(absPath)
	buf := bytes.NewBuffer(content)
	if err != nil {
		return nil, err
	}

	absPath = filepath.FromSlash(path.Join(StaticPath, name))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		return nil, err
	}

	absPath = filepath.FromSlash(path.Join(StaticPath, "Footer.html"))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), err
}
