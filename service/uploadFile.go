package service

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "method not supported ")
		return
	}
	err := r.ParseMultipartForm(10 << 20)
	fmt.Println("上传文件的路径前缀 ", CurrentPath)
	if err != nil {
		fmt.Fprintf(w, "上传的文件过大")
		return
	}

	file, header, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		fmt.Fprintf(w, "获取不到文件")
		return
	}
	log.Println("上传的文件信息： ", header.Filename, header.Size)
	buffer := make([]byte, 1024*1024)
	writePath := CurrentPath + header.Filename
	if isFileExisted(writePath) {
		fmt.Fprintf(w, "文件已经存在")
		return
	}

	WFile, err := os.OpenFile(writePath, os.O_CREATE|os.O_WRONLY, 0666)
	defer WFile.Close()
	if err != nil {
		log.Println("文件打开失败", err)
		return
	}
	for {
		n, err := file.Read(buffer)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			log.Println("上传文件失败")
			return
		}
		WFile.Write(buffer[:n])
	}
	fmt.Fprintf(w, "文件上传成功")
}

func isFileExisted(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
