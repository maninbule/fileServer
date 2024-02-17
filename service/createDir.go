package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func CreateDir(w http.ResponseWriter, r *http.Request) {
	log.Printf("enter CreateDir")
	if r.Method != http.MethodPost {
		fmt.Fprintf(w, "method not supported")
		return
	}
	r.ParseForm()
	dirName := r.Form.Get("dirName")
	curDir := r.Form.Get("curDir")
	path := curDir + dirName
	if createDirOnDisk(path) {
		fmt.Fprintf(w, "创建成功")
	}
}

func createDirOnDisk(path string) bool {
	err := os.Mkdir(path, 0666)
	if err != nil {
		log.Fatalf("创建文件夹失败 ", err)
		return false
	}
	return true
}
