package utils

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

func WriteHyperlinkListToRespon(w http.ResponseWriter, pathList []string, isDir []bool) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for i := 0; i < len(pathList); i++ {
		path := pathList[i]
		url := path
		if isDir[i] {
			url = string(filepath.Separator) + "Dir" + string(filepath.Separator) + url
		} else {
			url = string(filepath.Separator) + "file" + string(filepath.Separator) + url
		}
		road := strings.Split(path, "/")
		filename := road[len(road)-1]
		fmt.Fprintf(w, "<a href=\"%s\">%s</a><br/>", url, filename)
	}
}
