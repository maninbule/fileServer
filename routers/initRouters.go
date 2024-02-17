package routers

import (
	"github.com/fileServer/middleware"
	"github.com/fileServer/service"
	"github.com/fileServer/type/symbol"
	"log"
	"net/http"
)

func InitRouters() {
	fileServer := http.FileServer(http.Dir(symbol.StaticPrefix))
	prefixServer := http.StripPrefix(symbol.FilePrefix, fileServer)
	mux := http.NewServeMux()

	mux.Handle(symbol.FilePrefix, prefixServer)
	mux.HandleFunc("/index", service.FileListShow)
	mux.HandleFunc(symbol.DirPrefix, service.FileListShow)
	mux.Handle("/upload", middleware.LimitRequestBodyMiddleware(http.HandlerFunc(service.UploadFile)))
	//mux.HandleFunc("/upload", service.UploadFile)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatalf("服务启动失败 ", err)
	}
}
