package utils

import (
	"crypto/md5"
	"fmt"
	//"github.com/mitchellh/mapstructure"
	"crypto/sha1"
	"github.com/suboat/go-filepool/lib/go-uuid/uuid"
	"github.com/suboat/go-filepool/lib/mapstructure"
	"github.com/suboat/go-filepool/lib/mergo"
	"github.com/suboat/go-filepool/lib/scrypt"
	"io"
	"mime/multipart"
)

// md5 sum of string
func MD5(s string) (hash string) {
	h := md5.New()
	io.WriteString(h, s)
	hash = fmt.Sprintf("%x", h.Sum(nil))
	return
}

// update structure whith map
func UpdateStruct(input interface{}, source interface{}) (err error) {
	var (
		md      mapstructure.Metadata
		decoder *mapstructure.Decoder
	)
	config := &mapstructure.DecoderConfig{
		Metadata: &md,
		Result:   source,
	}
	if decoder, err = mapstructure.NewDecoder(config); err != nil {
		return
	}
	if err = decoder.Decode(input); err != nil {
		return
	}
	// 检查结构体的合法性
	if err = StrcutValid(source); err != nil {
		return
	}
	return
}

// 合并struct todo: 检查不允许合并的项
func MergeStruct(dst, src interface{}) error {
	return mergo.Merge(dst, src)
}

// scrypt string to key bit arrary
func ScryptString(s string, salt string) (dk []byte, err error) {
	dk, err = scrypt.Key([]byte(s), []byte(salt), 16384, 8, 1, 32)
	return
}

// 统一产生全站的uuid
func GenUuid() string {
	return uuid.New()
}

// 判断uuid是否合法
func UuidValid(s string) bool {
	if len(s) > 1 {
		return true
	} else {
		return false
	}
}

// multipart文件hash
func MultipartFileHash(file multipart.File) (hash string, err error) {
	h := sha1.New()
	if _, err = io.Copy(h, file); err != nil {
		return
	}
	hash = fmt.Sprintf("%x", h.Sum(nil))
	_, err = file.Seek(0, 0)
	return
}
