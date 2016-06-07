package upload

import (
	"encoding/json"
	"fmt"
	"github.com/suboat/go-filepool/config"
	"github.com/suboat/go-filepool/lib/resize"
	"github.com/suboat/go-filepool/utils"
	"github.com/suboat/sorm"
	"github.com/suboat/sorm/log"
	_ "golang.org/x/image/bmp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path"
	"time"
)

// 储存目录
var (
	StoreDir              string // 储存目录
	StoreDirThumbnail     string // 缩略图储存目录
	ResizeThumbnailWidth  uint   // 缩略图宽度
	ResizeThumbnailHeight uint   // 缩略图高度
)

type Picture struct {
	OriginName          string        `json:",omitempty"` // 原文件名
	OriginWidth         int           `json:",omitempty"` // 原文件宽度
	OriginHeight        int           `json:",omitempty"` // 原文件高度
	Size                int64         `json:",omitempty"` // 文件大小
	ContentType         string        `json:",omitempty"` // 文件类型
	StoreToken          orm.Accession `json:",omitempty"` // 储存key todo: 即ID
	DeleteToken         orm.Uid       `json:",omitempty"` // 删除key
	FetchToken          string        `json:",omitempty"` // 访问key
	FetchThumbnailToken string        `json:",omitempty"` // 访问缩略图key
	Sha1                string        `json:",omitempty"` // 文件哈希
	Sha1Thumbnail       string        `json:",omitempty"` // 缩略图哈希
	CreatTime           time.Time     `json:",omitempty"` // 上传时间
	DownloadTimes       int           `json:",omitempty"` // 下载次数
	Category            string        `json:",omitempty"` // 图片类型
	isInit              bool          `json:",omitempty"` // 是否已初始化
}

// json
func (p *Picture) ToJson() (s string) {
	s = "{}"
	if b, err := json.Marshal(p); err == nil {
		s = string(b)
	}
	return
}

// 初始化
func init() {
	var err error
	// 储存目录
	StoreDir = config.UploadPicStore
	StoreDirThumbnail = config.UploadPicStoreThumbnail
	if err = os.MkdirAll(StoreDir, os.ModePerm); err != nil {
		panic(err)
	}
	if err = os.MkdirAll(StoreDirThumbnail, os.ModePerm); err != nil {
		panic(err)
	}
	// 缩略图
	ResizeThumbnailWidth = 320
	ResizeThumbnailHeight = 0

	//// 普通索引
	//if err = models.ModelPic.EnsureIndex(orm.Index{
	//	"Key":      []string{"fetchtoken", "fetchthumbnailtoken", "category"},
	//	"Unique":   false,
	//	"DropDups": false, // important
	//}); err != nil {
	//	panic(err)
	//}

	//// Unique索引
	//if err = models.ModelPic.EnsureIndex(orm.Index{
	//	"Key":      []string{"storetoken"},
	//	"Unique":   true,
	//	"DropDups": true, // important
	//}); err != nil {
	//	panic(err)
	//}

	//config.Logger.Println("[DEBUG] models.picture init finish.")
}

// ***** 定义picture方法开始  *****
// 初始化数据
func (pic *Picture) Init() (err error) {
	if pic.isInit == false {
		pic.GenSavToken()
		pic.GenDelToken()
		pic.CreatTime = time.Now()

		pic.isInit = true
	}
	return
}

// 是否合法
func (pic *Picture) isValid() (err error) {
	return
}
func (pic *Picture) IsValid() (err error) {
	return pic.isValid()
}

// 产生删除token
func (pic *Picture) GenDelToken() (err error) {
	//pic.DeleteToken = models.ModelPic.NewUid()
	return
}

// 产生储存token
func (pic *Picture) GenSavToken() (err error) {
	pic.StoreToken = orm.NewAccession()
	return
}

// 保存到数据库
func (pic *Picture) Save() (err error) {
	//err = models.ModelPic.Objects().Create(pic)
	return
}

// 更新到数据库
func (pic *Picture) Update() (err error) {
	return
}

// 保存图片及缩略图
func (pic *Picture) SaveFileAndThumbnail(f multipart.File) (err error) {
	var (
		img         image.Image    // 图片
		img_c       image.Config   // 图片配置信息
		img_s       string         // 图片类型
		img_thm     image.Image    // 图片的缩略图
		img_thm_s   = "png"        // 缩略图类型
		img_thm_f   multipart.File // 缩略图文件句柄
		img_thm_sh1 string         // 缩略图的哈希
		out_img     *os.File       // 图片
		out_img_thm *os.File       // 缩略图
	)
	// 确认初始化
	if err = pic.Init(); err != nil {
		log.Error(err)
		return
	}
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
	pic.OriginWidth = img_c.Width
	pic.OriginHeight = img_c.Height
	// jpg后缀处理
	if img_s == "jpeg" {
		img_s = "jpg"
	}
	// 趋向转jpg缩略图
	if img_s == "jpg" {
		img_thm_s = "jpg"
	}
	pic.ContentType = img_s
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
	pic.Sha1Thumbnail = img_thm_sh1
	defer img_thm_f.Close()
	// 文件名
	pic.FetchThumbnailToken = fmt.Sprintf("%s.%s", pic.Sha1Thumbnail, img_thm_s)
	pic.FetchToken = fmt.Sprintf("%s.%s", pic.Sha1, pic.ContentType)
	p1 := path.Join(StoreDir, pic.FetchToken)
	//p2 := path.Join(StoreDirThumbnail, pic.FetchThumbnailToken)
	p2 := path.Join(StoreDirThumbnail, pic.FetchToken) // 缩略图与文件名相同
	log.Debug(p1, " - ", p2)
	// println("aa", "bb", p1, p2, img_c.Width, img_c.Height)
	if _, err = os.Stat(p1); err == nil { // 文件已存在，不做保存操作。这里认为图片和缩略图共同存在
		if _, err = os.Stat(p2); err == nil {
			return
		}
	}
	// 保存
	if out_img, err = os.Create(p1); err != nil {
		return
	}
	if out_img_thm, err = os.Create(p2); err != nil {
		return
	}
	defer func() {
		out_img.Close()
		out_img_thm.Close()
		if err != nil {
			_ = os.Remove(p1)
			_ = os.Remove(p2)
		}
	}()
	if _, err = io.Copy(out_img, f); err != nil {
		return
	}
	if _, err = io.Copy(out_img_thm, img_thm_f); err != nil {
		return
	}

	return
}

// 删除图片及缩略图
