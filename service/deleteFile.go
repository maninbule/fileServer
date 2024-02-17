package service

import (
	"fmt"
	"github.com/fileServer/utils"
	"log"
	"net/http"
	"os"
)

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	log.Printf("enter DeleteFile")
	if r.Method != "POST" {
		fmt.Fprintf(w, "method not supported")
		return
	}
	r.ParseForm()
	path := r.Form.Get("deletepath")
	log.Printf(path)
	if deleteFileFromDisk(path) {
		fmt.Fprintf(w, "删除成功")
	} else {
		fmt.Fprintf(w, "删除失败")
	}
}

func deleteFileFromDisk(path string) bool {
	if utils.IsFileExisted(path) {
		err := os.RemoveAll(path)
		if err != nil {
			log.Printf("删除文件失败 ", err)
			return false
		}
	} else {
		return false
	}
	return true
}
