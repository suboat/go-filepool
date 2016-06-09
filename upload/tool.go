package upload

import (
	"github.com/suboat/go-filepool/config"
	"github.com/suboat/go-filepool/lib/resize"
	"github.com/suboat/go-filepool/utils"
	"github.com/suboat/sorm/log"
	"mime/multipart"
	"path"
	"path/filepath"

	_ "golang.org/x/image/bmp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"io"
	"os"
)

// 检查没有缩略图的图片文件,并生成缩略图
func ToolFixThumbnail() (err error) {
	log.Debug("Walk ", config.UploadPicStore)
	filepath.Walk(config.UploadPicStore, toolFixThumbnail)
	return
}

func toolFixThumbnail(p string, fi os.FileInfo, er error) (err error) {
	if fi.IsDir() {
		return
	}
	// 缩略图路径
	var thumbPath = path.Join(config.UploadPicStoreThumbnail, fi.Name())

	// 缩略图已存在
	if _, _err := os.Stat(thumbPath); _err == nil {
		return
	}

	// 如果缩略图不存在则尝试生成
	log.Debug("[NOT EXIST] ", thumbPath)
	fp, _ := os.Open(p)
	defer fp.Close()
	if er = ToolThumbnail(fp, thumbPath); er != nil {
		log.Error("[FAIL] ", p, " ", er)
		return
	}
	log.Info("[FIX THUMB] ", thumbPath)

	return
}

// 将图片转缩略图
func ToolThumbnail(f multipart.File, p2 string) (err error) {
	var (
		img         image.Image    // 图片
		img_c       image.Config   // 图片配置信息
		img_s       string         // 图片类型
		img_thm     image.Image    // 图片的缩略图
		img_thm_s   = "png"        // 缩略图类型
		img_thm_f   multipart.File // 缩略图文件句柄
		img_thm_sh1 string         // 缩略图的哈希
		out_img_thm *os.File       // 缩略图
	)

	// 解析图片
	if img, img_s, err = image.Decode(f); err != nil {
		log.Error(err)
		return
	}
	if _, err = f.Seek(0, 0); err != nil {
		log.Error(err)
		return
	}
	// 解析图片config信息,宽高
	if img_c, _, err = image.DecodeConfig(f); err != nil {
		log.Error(err)
		return
	}
	// jpg后缀处理
	if img_s == "jpeg" {
		img_s = "jpg"
	}
	// 趋向转jpg缩略图
	if img_s == "jpg" {
		img_thm_s = "jpg"
	}
	if _, err = f.Seek(0, 0); err != nil {
		log.Error(err)
		return
	}
	// 转缩略图
	img_thm = resize.Resize(ResizeThumbnailWidth, ResizeThumbnailHeight, img, resize.NearestNeighbor)
	if img_thm_f, img_thm_sh1, err = utils.ImageToiMutiFile(img_thm, img_thm_s); err != nil {
		log.Error(err)
		return
	}
	defer img_thm_f.Close()
	// 保存
	if out_img_thm, err = os.Create(p2); err != nil {
		return
	}
	defer func() {
		out_img_thm.Close()
		if err != nil {
			_ = os.Remove(p2)
		}
	}()
	if _, err = io.Copy(out_img_thm, img_thm_f); err != nil {
		return
	}
	// debug
	log.Debug("thunb hash:", img_thm_sh1, " org:", img_c.Width, "x", img_c.Height)
	return
}
