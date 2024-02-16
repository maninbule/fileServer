package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	fileServer := http.FileServer(http.Dir("/static"))
	http.Handle("/", fileServer)

	http.HandleFunc("/index", fileListShow)

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
	if path == "/index" {
		path = "static"
	}
	log.Println(path)
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
		for _, value := range dir {
			nextPath := filepath.Join(path, string(filepath.Separator), value.Name())
			if value.IsDir() {
				nextPath += "/"
			}
			fileList = append(fileList, nextPath)
		}
		for _, p := range fileList {
			fmt.Fprintln(w, p)
		}
	}
}
