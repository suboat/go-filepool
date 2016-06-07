package filepool

import (
	"mime/multipart"
	"net/http"
)

// 判断文件类型，检测前512字节数据，如果不能识别，返回"application/octet-stream"
func DetectContentType(f multipart.File) (t string, err error) {
	p := make([]byte, 512, 512)
	_, err = f.Read(p)
	if err == nil {
		t = http.DetectContentType(p)
		_, err = f.Seek(0, 0)
	}
	return
}

// 判断是否图片文件
func DetectImageType(f multipart.File) (t string, err error) {
	if t, err = DetectContentType(f); err != nil {
		return
	}
	err = IsImageType(t)
	return
}

func IsImageType(t string) (err error) {
	switch t {
	case "image/jpeg", "image/png", "image/gif", "image/bmp":
		err = nil
	default:
		println(t)
		err = ErrImageType
	}
	return
}
