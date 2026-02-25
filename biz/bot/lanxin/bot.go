package lanxin
package lanxin

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	SecurityNone    = "none"
	SecuritySign    = "sign"
	SecurityKeyword = "keyword"
)

// Robot 蓝信机器人
type Robot struct {
	WebHookUrl   string
	Client       *http.Client
	SecurityType string
	Secret       string
	Keywords     string
}

// Resp 蓝信API响应
type Resp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// New 创建蓝信机器人实例（传入完整的webhook URL）
func New(webhookUrl string) *Robot {
	return &Robot{
		WebHookUrl:   webhookUrl,
		Client:       &http.Client{Timeout: 5 * time.Second},
		SecurityType: SecurityNone,
	}
}

// AddSign 设置签名模式
func (r *Robot) AddSign(secret string) *Robot {
	if secret != "" {
		r.SecurityType = SecuritySign
		r.Secret = secret
	}
	return r
}

// AddKeyword 设置关键字模式
func (r *Robot) AddKeyword(keyword string) *Robot {
	if keyword != "" {
		r.SecurityType = SecurityKeyword
		r.Keywords = keyword
	}
	return r
}

// Send 发送消息
func (r *Robot) Send(message interface{}) (bool, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	return r.SendRaw(b)
}

// SendRaw 发送原始JSON消息
func (r *Robot) SendRaw(msgBytes []byte) (bool, error) {
	requestUrl := r.WebHookUrl

	// 签名模式：在URL中追加timestamp和sign参数
	if r.SecurityType == SecuritySign && r.Secret != "" {
		timestamp := time.Now().UnixNano() / 1e6
		sign := r.sign(timestamp, r.Secret)
		requestUrl = fmt.Sprintf("%s&timestamp=%d&sign=%s", requestUrl, timestamp, url.QueryEscape(sign))
	}

	// 关键字模式：在消息内容中注入关键字
	if r.SecurityType == SecurityKeyword && r.Keywords != "" {
		msgBytes = r.injectKeyword(msgBytes)
	}

	resp, err := r.Client.Post(requestUrl, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	ret := &Resp{}
	if err := json.Unmarshal(body, ret); err != nil {
		return false, err
	}

	if ret.ErrCode == 0 {
		return true, nil
	}
	return false, errors.New(ret.ErrMsg)
}

// sign 生成HMAC-SHA256签名
func (r *Robot) sign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// injectKeyword 在消息内容中注入关键字
func (r *Robot) injectKeyword(msgBytes []byte) []byte {
	var raw map[string]interface{}
	if err := json.Unmarshal(msgBytes, &raw); err != nil {
		return msgBytes
	}

	injected := false
	// text 消息
	if text, ok := raw["text"].(map[string]interface{}); ok {
		if content, ok := text["content"].(string); ok {
			text["content"] = r.Keywords + "\n" + content
			injected = true
		}
	}
	// markdown 消息
	if !injected {
		if md, ok := raw["markdown"].(map[string]interface{}); ok {
			if text, ok := md["text"].(string); ok {
				md["text"] = r.Keywords + "\n" + text
				injected = true
			}
		}
	}

	newBytes, err := json.Marshal(raw)
	if err != nil {
		return msgBytes
	}
	return newBytes
}

// CheckMessage 检查消息是否为合法的蓝信消息格式
func (r *Robot) CheckMessage(msg string) bool {
	if len(msg) == 0 {
		return false
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &raw); err != nil {
		return false
	}
	msgType, ok := raw["msgtype"].(string)
	if !ok {
		return false
	}
	switch msgType {
	case "text", "markdown":
		return true
	default:
		return false
	}
}
