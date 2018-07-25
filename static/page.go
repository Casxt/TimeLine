package static

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Casxt/TimeLine/config"
)

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

	return 200, buf.Bytes(), err
}

//GetFile Get static file
func GetFile(names ...string) (status int, File []byte, err error) {
	var absPath string
	var content []byte
	//res http.ResponseWriter,
	absPath = filepath.FromSlash(path.Join(config.ProjectPath, path.Join(names...)))
	content, err = ioutil.ReadFile(absPath)

	switch {
	case err == nil:
	case strings.HasSuffix(err.Error(), "The system cannot find the file specified."):
		log.Println("GetFile", absPath, "Not Found.")
		return 404, nil, errors.New("File Not Found")
	default:
		log.Println("GetFile", err.Error())
		return 500, nil, err
	}
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
