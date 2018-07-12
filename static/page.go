package static

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Casxt/TimeLine/config"
)

//ProjectPath is path that save the static file
//var ProjectPath = "C:\\Users\\Surface\\go\\src\\github.com\\Casxt\\TimeLine"

//GetPage return the builted html
func GetPage(names ...string) (status int, Page []byte, err error) {
	var absPath string
	var content []byte
	var buf *bytes.Buffer

	absPath = filepath.FromSlash(path.Join(config.ProjectPath, "static", "Header.html"))
	content, err = ioutil.ReadFile(absPath)
	buf = bytes.NewBuffer(content)
	if err != nil {
		absPath = filepath.FromSlash(path.Join(config.ProjectPath, "static", "404.html"))
		content, err = ioutil.ReadFile(absPath)
		buf = bytes.NewBuffer(content)
		return 404, buf.Bytes(), err
	}

	absPath = filepath.FromSlash(path.Join(config.ProjectPath, path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		absPath = filepath.FromSlash(path.Join(config.ProjectPath, "static", "404.html"))
		content, err = ioutil.ReadFile(absPath)
		buf = bytes.NewBuffer(content)
		return 404, buf.Bytes(), err
	}

	absPath = filepath.FromSlash(path.Join(config.ProjectPath, "static", "Footer.html"))
	content, err = ioutil.ReadFile(absPath)
	buf.Write(content)
	if err != nil {
		absPath = filepath.FromSlash(path.Join(config.ProjectPath, "static", "404.html"))
		content, err = ioutil.ReadFile(absPath)
		buf = bytes.NewBuffer(content)
		return 404, buf.Bytes(), err
	}

	return 404, buf.Bytes(), err
}

//GetFile Get static file
func GetFile(names ...string) (status int, File []byte, err error) {
	var absPath string
	var content []byte
	//res http.ResponseWriter,
	absPath = filepath.FromSlash(path.Join(config.ProjectPath, path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)

	if err != nil {
		log.Println(err.Error())
		return 500, nil, err
	}
	//log.Println(mime.TypeByExtension(path.Ext(names[len(names)-1])))
	//log.Println(path.Ext(names[len(names)-1]))
	//res.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(names[len(names)-1])))
	return 200, content, err
}

//SaveFile save static file
func SaveFile(bytes []byte, names ...string) (err error) {
	var absPath string
	absPath = filepath.FromSlash(path.Join(config.ProjectPath, path.Join(names...)))
	//os.ModeAppend in there will not append, but clear and write
	//if file not exist, will create file
	err = ioutil.WriteFile(absPath, bytes, os.ModeAppend)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
