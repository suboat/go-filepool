package filepool

import (
	"fmt"
	"github.com/suboat/go-filepool/utils"
	"github.com/suboat/sorm/log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type UploadFileRequire struct {
	MaxSize      int64  // 上传文件大小限制
	FormName     string // 表格中储存file的字段名
	RequireImage bool   // 是否要求图片类型
	Rename       string // 将文件名视为
}

// 获取到的文件
type UploadFileMode struct {
	File        multipart.File // 文件
	FileName    string         // 文件名
	Size        int64          // 文件大小
	Sha1        string         // sha1值
	ContentType string         // 文件类型
}

func (f *UploadFileMode) Close() {
	if f.File != nil {
		f.File.Close()
	}
}

// 获取文件大小的接口
type checkSize interface {
	Size() int64
}

// 获取文件信息的接口
type checkStat interface {
	Stat() (os.FileInfo, error)
}

func UploadFileOne(rw http.ResponseWriter, req *http.Request, meta *UploadFileRequire) (*UploadFileMode, error) {
	var (
		file   multipart.File
		header *multipart.FileHeader
		res    = new(UploadFileMode)
		err    error
	)

	// 限制用于解析file的内存大小，超出部分保存在硬盘
	req.ParseMultipartForm(32 << 20)
	file, header, err = req.FormFile(meta.FormName)
	if err != nil {
		rw.WriteHeader(500)
		fmt.Fprint(rw, err.Error())
		log.Error(err)
		return nil, err
	}

	// 文件保存在内存，sectionReader
	// 文件保存在硬盘，os.FileInfo
	if s, ok := file.(checkSize); ok {
		res.Size = s.Size()
	} else if s, ok := file.(checkStat); ok {
		_info, _ := s.Stat()
		res.Size = _info.Size()
	}
	res.FileName = header.Filename
	res.File = file

	// 文件大小限制
	if meta.MaxSize > 0 && (res.Size > meta.MaxSize) {
		err = ErrUploadFileSize
		return nil, err
	}

	// 文件类型限制
	if res.ContentType, err = DetectContentType(res.File); err != nil {
		return nil, err
	}
	log.Debug("file type", res.ContentType)
	if meta.RequireImage == true {
		if err = IsImageType(res.ContentType); err != nil {
			return nil, err
		}
	}

	// sha1 值
	res.Sha1, err = utils.MultipartFileHash(res.File)

	return res, err
}

// 普通文件转为UploadFileMode类型
func NormalToUploadFile(p string, meta *UploadFileRequire) (res *UploadFileMode, err error) {
	var (
		f  *os.File
		fi os.FileInfo
	)
	if fi, err = os.Stat(p); err != nil {
		return nil, err
	}
	if f, err = os.Open(p); err != nil {
		return nil, err
	}
	res = new(UploadFileMode)
	res.File = multipart.File(f)
	res.FileName = path.Base(p)
	res.Size = fi.Size()

	// 将文件名视为
	if meta != nil && len(meta.Rename) > 0 {
		res.FileName = meta.Rename
	}

	// 文件大小限制
	if meta != nil && meta.MaxSize > 0 && (res.Size > meta.MaxSize) {
		err = ErrUploadFileSize
		return nil, err
	}

	// 文件类型
	if res.ContentType, err = DetectContentType(res.File); err != nil {
		return nil, err
	}

	// 文件类型限制
	if meta != nil && meta.RequireImage == true {
		if err = IsImageType(res.ContentType); err != nil {
			return nil, err
		}
	}

	// sha1 值
	res.Sha1, err = utils.MultipartFileHash(res.File)

	return
}
