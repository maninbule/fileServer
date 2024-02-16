package main

import (
	"fmt"
	"github.com/fileServer/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const StaticPrefix string = "static" + string(filepath.Separator) + "file" + string(filepath.Separator)
const FilePrefix string = "/file/"
const DirPrefix string = "/Dir/"

func main() {
	fileServer := http.FileServer(http.Dir(StaticPrefix))
	prefixServer := http.StripPrefix(FilePrefix, fileServer)
	http.Handle(FilePrefix, prefixServer)
	http.HandleFunc("/index", fileListShow)
	http.HandleFunc(DirPrefix, fileListShow)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("服务启动失败 ", err)
	}
}

func fileListShow(w http.ResponseWriter, r *http.Request) {
	log.Println("enter fileListShow")
	if r.Method != "GET" {
		http.Error(w, "method is not supported", http.StatusBadRequest)
		return
	}
	path := r.URL.Path
	log.Println(path)
	if path == "/index" {
		path = StaticPrefix
	} else if strings.HasPrefix(path, DirPrefix) {
		path = strings.TrimPrefix(path, DirPrefix)
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
			nextPath := filepath.Join(path, string(filepath.Separator), subFile.Name())
			if subFile.IsDir() {
				nextPath += string(filepath.Separator)
			}
			if subFile.IsDir() {
				isDir = append(isDir, true)
			} else {
				nextPath = strings.TrimPrefix(nextPath, StaticPrefix)
				isDir = append(isDir, false)
			}
			fileList = append(fileList, nextPath)
		}
		fmt.Println(fileList, isDir)
		utils.WriteHyperlinkListToRespon(w, fileList, isDir)
	}
}
