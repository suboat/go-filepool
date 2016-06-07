package config

import (
	"os"
	"path"
	"path/filepath"
)

const (
	UploadFileCateDefault = "unknow"  // 上传文件默认分类：未分类
	PicCateDefault        = "default" // 上传图片文件默认分类：未分类
	PicCateItem           = "item"    //
)

var (
	BaseDir                 string // 程序目录
	DownloadDir             string // 下载目录
	UploadPicStore          string // 上传图片保存的正式目录
	UploadPicStoreThumbnail string // 上传图片保存的缩略图目录
)

func init() {
	BaseDir = path.Join(path.Dir(os.Args[0]), "./") // 运行目录
	Init()
}

func Init() {
	BaseDir, _ = filepath.Abs(BaseDir)                                         // 绝对路径
	DownloadDir = path.Join(BaseDir, "upload")                                 // 下载目录
	UploadPicStore = path.Join(BaseDir, "upload", "pictrue")                   // 上传图片保存的正式目录
	UploadPicStoreThumbnail = path.Join(BaseDir, "upload", "pictruethumbnail") // 上传图片保存的缩略图目录
}
