package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/suboat/go-filepool/config"
	"github.com/suboat/go-filepool/upload"
	"github.com/suboat/sorm/log"

	"github.com/go-ini/ini"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

var (
	CfgMap *mainCfg = nil
)

// 简单配置文件
type mainCfg struct {
	FilePath    string `ini:"-"`
	Address     string
	UploadUrl   string
	DownloadUrl string
}

func newMainCfg(s *mainCfg) (d *mainCfg, err error) {
	d = &mainCfg{
		Address:     "0.0.0.0:8091",
		UploadUrl:   "/upload",
		DownloadUrl: "/download/",
	}
	d.FilePath, _ = filepath.Abs(path.Join(path.Dir(os.Args[0]), "./"))
	d.FilePath = path.Join(d.FilePath, "config.ini")
	return
}

func main() {

	// log
	//log.SetLevel(logrus.DebugLevel)
	log.SetLevel(logrus.InfoLevel)

	// 配置文件路径 及默认设置
	var (
		cfg *ini.File
		err error
	)
	CfgMap, _ = newMainCfg(nil)
	if cfg, err = ini.Load(CfgMap.FilePath); err != nil {
		//panic(err)
		if cfg, err = ini.LooseLoad(CfgMap.FilePath); err != nil {
			panic(err)
		}
		if err = ini.ReflectFrom(cfg, CfgMap); err != nil {
			panic(err)
		}
		if err = cfg.SaveTo(CfgMap.FilePath); err != nil {
			return
		}
	} else {
		if err = cfg.MapTo(CfgMap); err != nil {
			panic(err)
		}
	}

	var (
		h = &upload.UploadHandler{}
	)
	h.FormName = "file"          // 文件名
	h.RequireImage = true        // 要求是图片
	h.MaxSize = 20 * 1024 * 1024 // 大小限制 20MB

	// 上传
	http.Handle(CfgMap.UploadUrl, h)

	// 下载
	log.Info("FileDir:", config.DownloadDir, " DownloadUrl:", CfgMap.DownloadUrl)
	http.Handle(CfgMap.DownloadUrl, config.DirNotList(
		http.StripPrefix(CfgMap.DownloadUrl, http.FileServer(http.Dir(config.DownloadDir)))))

	// 检查缺失的缩略图
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		if err := upload.ToolFixThumbnail(); err != nil {
			log.Error(err)
		}
	}()

	// run
	log.Info("ListenAndServe ", CfgMap.Address)
	http.ListenAndServe(CfgMap.Address, nil)
}
