package upload

import (
	"encoding/json"
	"fmt"
	"github.com/suboat/go-filepool"
	"github.com/suboat/sorm/log"
	"net/http"
)

type UploadHandler struct {
	filepool.UploadFileRequire
	Category string
}

type ErrorResp struct {
	Error    error  `json:"-"`
	ErrorStr string `json:"error,omitempty"`
}

// json
func (e *ErrorResp) ToJson() (s string) {
	s = "{}"
	if e != nil && e.Error != nil {
		e.ErrorStr = e.Error.Error()
	}
	if b, err := json.Marshal(e); err == nil {
		s = string(b)
	}
	return
}

func (h *UploadHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var (
		uploadRequire = &filepool.UploadFileRequire{
			MaxSize:      h.MaxSize,
			FormName:     h.FormName,
			RequireImage: h.RequireImage,
		}
		p = &Picture{Category: h.Category}
		//r    = response.ResponseNew()
		file *filepool.UploadFileMode
	)

	var (
		resp *ErrorResp = new(ErrorResp)
	)

	defer func() {
		if resp.Error != nil {
			fmt.Fprint(rw, resp.ToJson())
		}
	}()

	// 解析上传文件
	if file, resp.Error = filepool.UploadFileOne(rw, req, uploadRequire); resp.Error != nil {
		// err
		log.Error(resp.Error)
		return
	}
	defer file.Close()
	p.OriginName = file.FileName
	p.Size = file.Size
	p.Sha1 = file.Sha1

	// 保存原图和缩略图
	if resp.Error = p.SaveFileAndThumbnail(file.File); resp.Error != nil {
		//r.ErrorSet(err)
		log.Error(resp.Error)
		return
	}

	//// 保存到数据库
	//if err = p.Save(); err != nil {
	//	return
	//}

	// 返回结果
	fmt.Fprint(rw, p.ToJson())

	return
}

//func NewHandler(h *UploadHandler) response.RestHandler {
//	if h == nil {
//		h = &UploadHandler{}
//	}
//	// 补全参数
//	if h.MaxSize == 0 {
//		h.MaxSize = 10 * 1024 * 1024 // 10MB上传限制
//	}
//	if h.FormName == "" {
//		h.FormName = "file" // form file name
//	}
//	if h.Category == "" {
//		h.Category = "default"
//	}
//	return h
//}
//
//func NewHandlerPic(h *UploadHandler) response.RestHandler {
//	if h == nil {
//		h = &UploadHandler{}
//	}
//	h.RequireImage = true // 要求是图片
//	// 默认图片分类
//	if h.Category == "" {
//		h.Category = "default"
//	}
//	return NewHandler(h)
//}
