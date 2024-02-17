package service

import (
	"fmt"
	"github.com/fileServer/type/symbol"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

var CurrentPath string = ""

func FileListShow(w http.ResponseWriter, r *http.Request) {
	log.Println("enter fileListShow")
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusBadRequest)
		return
	}
	path := r.URL.Path
	log.Println(path)
	if path == "/index" {
		path = symbol.StaticPrefix
	} else if strings.HasPrefix(path, symbol.DirPrefix) {
		path = strings.TrimPrefix(path, symbol.DirPrefix)
	}
	CurrentPath = path
	file, err := os.Open(path)
	if err != nil {
		log.Println("os.open error ", err)
		http.Error(w, "404 not found", http.StatusBadRequest)
		return
	}
	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		return
	}
	if stat.IsDir() {
		dir, err := file.ReadDir(-1)
		if err != nil {
			http.Error(w, "服务器内部错误", http.StatusInternalServerError)
			return
		}
		fileList := make([]string, 0)
		isDir := make([]bool, 0)
		for _, subFile := range dir {
			nextPath := path + subFile.Name()
			if subFile.IsDir() {
				isDir = append(isDir, true)
			} else {
				nextPath = strings.TrimPrefix(nextPath, symbol.StaticPrefix)
				isDir = append(isDir, false)
			}
			fileList = append(fileList, nextPath)
		}
		fmt.Println(fileList, isDir)
		//writeHyperlinkListToRespon(w, fileList, isDir)
		pathList, filename := pathTohtml(fileList, isDir)
		useTemplate(w, pathList, filename, CurrentPath)
	}
}

func writeHyperlinkListToRespon(w http.ResponseWriter, pathList []string, isDir []bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>Upload File</title>\n</head>\n<body>")
	for i := 0; i < len(pathList); i++ {
		path := pathList[i]
		url := path
		road := strings.Split(path, "/")
		filename := road[len(road)-1]
		if isDir[i] {
			url = "/" + "Dir" + "/" + url + "/"
			filename += "/"
		} else {
			url = "/" + "file" + "/" + url
		}

		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br/>", url, filename)
	}

	upload := "<form action=\"/upload\" method=\"post\" enctype=\"multipart/form-data\">\n    <input type=\"file\" name=\"file\"><br>\n    <input type=\"submit\" value=\"Upload\"> <br>\n    <input type=\"hidden\" name=\"url\" value=\"{{ .URL }}\">\n</form>"
	fmt.Fprintf(w, upload)
	fmt.Fprintf(w, "</body>\n</html>")
}

func pathTohtml(pathList []string, isDir []bool) ([]string, []string) {
	urls := make([]string, 0)
	name := make([]string, 0)
	for i := 0; i < len(pathList); i++ {
		path := pathList[i]
		url := path
		road := strings.Split(path, "/")
		filename := road[len(road)-1]
		if isDir[i] {
			url = "/" + "Dir" + "/" + url + "/"
			filename += "/"
		} else {
			url = "/" + "file" + "/" + url
		}
		urls = append(urls, url)
		name = append(name, filename)
	}
	return urls, name
}

func useTemplate(w http.ResponseWriter, urls []string, filename []string, curDir string) {
	tmpl, err := template.ParseFiles("static/template/index.html")
	if err != nil {
		log.Fatalf("模板引擎解析失败")
		return
	}
	data := struct {
		Urls     []string
		Filename []string
		CurDir   string
	}{
		Urls:     urls,
		Filename: filename,
		CurDir:   curDir,
	}
	fmt.Print(curDir)
	tmpl.Execute(w, data)
}
