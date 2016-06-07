package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/suboat/go-filepool/config"
	"github.com/suboat/go-filepool/upload"
	"github.com/suboat/sorm/log"

	"net/http"
)

func main() {

	// log
	log.SetLevel(logrus.DebugLevel)

	var (
		h       = &upload.UploadHandler{}
		address = "0.0.0.0:8091"
	)
	h.FormName = "file"          // 文件名
	h.RequireImage = true        // 要求是图片
	h.MaxSize = 20 * 1024 * 1024 // 大小限制 20MB

	// 上传
	http.Handle("/upload", h)

	// 下载
	log.Debug("DownloadDir ", config.DownloadDir)
	http.Handle("/download/", config.DirNotList(http.StripPrefix("/download/", http.FileServer(http.Dir(config.DownloadDir)))))

	// run
	log.Debug("ListenAndServe ", address)
	http.ListenAndServe(address, nil)
}
