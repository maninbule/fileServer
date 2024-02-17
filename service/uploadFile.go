package service

import (
	"fmt"
	"github.com/fileServer/utils"
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
	const maxUploadSize = 22<<20 + 512 // 12.5MB
	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		fmt.Fprintf(w, "解析表出错")
		return
	}
	fmt.Println("上传文件的路径前缀 ", CurrentPath)
	file, header, err := r.FormFile("file")
	pathDir := r.Form.Get("saveDir")
	defer file.Close()
	if err != nil {
		fmt.Fprintf(w, "获取不到文件")
		return
	}
	log.Println("上传的文件信息： ", header.Filename, header.Size)
	buffer := make([]byte, 1024*1024)
	writePath := pathDir + header.Filename
	if utils.IsFileExisted(writePath) {
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
