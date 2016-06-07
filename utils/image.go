package utils

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
)

// 图片对象转文件对象
// 返回multipart.File类型和sha1值
func ImageToiMutiFile(img image.Image, t string) (f multipart.File, hash string, err error) {
	var (
		buf = new(bytes.Buffer)
		b   []byte
	)
	// 转换
	switch t {
	case "jpeg":
		err = jpeg.Encode(buf, img, nil)
	case "jpg":
		err = jpeg.Encode(buf, img, nil)
	case "png":
		err = png.Encode(buf, img)
	case "gif":
		err = gif.Encode(buf, img, nil)
	}
	if err != nil {
		return
	}
	b = buf.Bytes()
	f = sectionReadCloser{io.NewSectionReader(bytes.NewReader(b), 0, int64(len(b)))}

	// sha1计算
	hash, err = MultipartFileHash(f)
	return
}

// helper types to turn a []byte into a File

type sectionReadCloser struct {
	*io.SectionReader
}

func (rc sectionReadCloser) Close() error {
	return nil
}
