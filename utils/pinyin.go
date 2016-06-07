package utils

import (
	"github.com/mozillazg/go-pinyin"
)

var punctuationMap = map[string]string{
	"（": "(",
	"）": ")",
	"、": ",",
	"，": ",",
	"－": "-",
	"！": "!",
	"？": "?",
}

func Pinyin(s string) (ret string) {
	a := pinyin.NewArgs()
	for _, r := range []rune(s) {
		if p, ok := punctuationMap[string(r)]; ok == true {
			ret += p
		} else {
			c := pinyin.SinglePinyin(r, a)
			if len(c) == 0 {
				ret += string(r)
			} else {
				ret += string(c[0])
			}
		}
	}
	return
}
