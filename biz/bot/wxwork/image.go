package wxwork

import (
	"crypto/md5"
	"encoding/hex"
)

type imageMessage struct {
	message
	Image Image `json:"image"`
}

type Image struct {
	Base64 string `json:"base64"`
	MD5    string `json:"md5"`
}

func NewImage() *Image {
	return &Image{}
}

func (i *Image) SetBase64(s string) *Image {
	i.Base64 = s
	h := md5.Sum([]byte(s))
	i.MD5 = hex.EncodeToString(h[:])
	return i
}
