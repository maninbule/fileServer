package service

import (
	"fmt"
	"github.com/fileServer/type/symbol"
	"log"
	"net/http"
	"os"
	"strings"
)

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
				nextPath += "/"
			}
			if subFile.IsDir() {
				isDir = append(isDir, true)
			} else {
				nextPath = strings.TrimPrefix(nextPath, symbol.StaticPrefix)
				isDir = append(isDir, false)
			}
			fileList = append(fileList, nextPath)
		}
		fmt.Println(fileList, isDir)
		writeHyperlinkListToRespon(w, fileList, isDir)
	}
}

func writeHyperlinkListToRespon(w http.ResponseWriter, pathList []string, isDir []bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for i := 0; i < len(pathList); i++ {
		path := pathList[i]
		url := path
		if isDir[i] {
			url = "/" + "Dir" + "/" + url
		} else {
			url = "/" + "file" + "/" + url
		}
		road := strings.Split(path, "/")
		filename := road[len(road)-1]
		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br/>", url, filename)
	}
}
