package routers

import (
	"github.com/fileServer/service"
	"github.com/fileServer/type/symbol"
	"log"
	"net/http"
)

func InitRouters() {
	fileServer := http.FileServer(http.Dir(symbol.StaticPrefix))
	prefixServer := http.StripPrefix(symbol.FilePrefix, fileServer)
	http.Handle(symbol.FilePrefix, prefixServer)
	http.HandleFunc("/index", service.FileListShow)
	http.HandleFunc(symbol.DirPrefix, service.FileListShow)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("服务启动失败 ", err)
	}
}
