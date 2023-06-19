package wxwork

import "github.com/yi-nology/common/utils/xhash"

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
	i.MD5 = xhash.Md5Hex([]byte(s))
	return i
}
