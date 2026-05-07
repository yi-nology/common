package xhash

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/spaolacci/murmur3"
)

// Hash returns the hash value of data.
func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

// MM3Md5 returns the md5 bytes of data.
func MM3Md5(data []byte) []byte {
	digest := md5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

// Md5Hex returns the md5 hex string of data.
func Md5Hex(data []byte) string {
	return fmt.Sprintf("%x", MM3Md5(data))
}

// 对字符串进行sha1 计算
func Sha1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

// 对数据进行md5计算
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 对数据进行 base64 编码
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// 对数据进行 base64 解码
func Base64Decode(s string) (string, error) {
	rs, err := base64.StdEncoding.DecodeString(s)
	return string(rs), err
}
