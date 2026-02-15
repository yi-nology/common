package lark

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const (
	defaultWebHookUrlTemplate = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"

	SecurityNone    = "none"
	SecuritySign    = "sign"
	SecurityKeyword = "keyword"
)

// LarkBot 飞书机器人
type LarkBot struct {
	Key          string
	WebHookUrl   string
	Client       *http.Client
	SecurityType string
	Secret       string
	Keywords     string
}

// Resp 飞书API响应
type Resp struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	Code          int    `json:"code"`
	Data          struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

// New 创建飞书机器人实例
func New(botKey string) *LarkBot {
	return &LarkBot{
		Key:          botKey,
		WebHookUrl:   fmt.Sprintf(defaultWebHookUrlTemplate, botKey),
		Client:       &http.Client{Timeout: 5 * time.Second},
		SecurityType: SecurityNone,
	}
}

// AddSign 设置签名模式
func (l *LarkBot) AddSign(secret string) *LarkBot {
	if secret != "" {
		l.SecurityType = SecuritySign
		l.Secret = secret
	}
	return l
}

// AddKeyword 设置关键字模式
func (l *LarkBot) AddKeyword(keyword string) *LarkBot {
	if keyword != "" {
		l.SecurityType = SecurityKeyword
		l.Keywords = keyword
	}
	return l
}

// Send 发送消息
func (l *LarkBot) Send(msg interface{}) (bool, error) {
	msgBytes, err := marshalMessage(msg)
	if err != nil {
		return false, err
	}

	// 签名模式：将 timestamp 和 sign 注入到 JSON body 中
	if l.SecurityType == SecuritySign && l.Secret != "" {
		msgBytes, err = l.injectSign(msgBytes)
		if err != nil {
			return false, err
		}
	}

	// 关键字模式：在消息内容中注入关键字
	if l.SecurityType == SecurityKeyword && l.Keywords != "" {
		msgBytes = l.injectKeyword(msgBytes)
	}

	return l.SendRaw(msgBytes)
}

// SendRaw 发送原始JSON消息
func (l *LarkBot) SendRaw(msgBytes []byte) (bool, error) {
	req, err := http.NewRequest(http.MethodPost, l.WebHookUrl, bytes.NewBuffer(msgBytes))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	r := Resp{}
	if err = json.Unmarshal(body, &r); err != nil {
		return false, err
	}
	if r.Code != 0 {
		return false, fmt.Errorf("send message failed, code: %d, msg: %s", r.Code, r.Msg)
	}

	return true, nil
}

// sign 生成飞书签名（timestamp + "\n" + secret -> HMAC-SHA256 -> Base64）
func (l *LarkBot) sign(timestamp int64, secret string) string {
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(stringToSign))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// injectSign 将签名信息注入到消息JSON body中
func (l *LarkBot) injectSign(msgBytes []byte) ([]byte, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(msgBytes, &raw); err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	raw["timestamp"] = strconv.FormatInt(timestamp, 10)
	raw["sign"] = l.sign(timestamp, l.Secret)

	return json.Marshal(raw)
}

// injectKeyword 在消息内容中注入关键字
func (l *LarkBot) injectKeyword(msgBytes []byte) []byte {
	var raw map[string]interface{}
	if err := json.Unmarshal(msgBytes, &raw); err != nil {
		return msgBytes
	}

	injected := false
	// text 消息: content.text
	if content, ok := raw["content"].(map[string]interface{}); ok {
		if text, ok := content["text"].(string); ok {
			content["text"] = l.Keywords + "\n" + text
			injected = true
		}
		// post 消息: content.post.zh_cn.title
		if !injected {
			if post, ok := content["post"].(map[string]interface{}); ok {
				if zhCn, ok := post["zh_cn"].(map[string]interface{}); ok {
					if title, ok := zhCn["title"].(string); ok {
						zhCn["title"] = l.Keywords + " " + title
						injected = true
					}
				}
			}
		}
	}

	newBytes, err := json.Marshal(raw)
	if err != nil {
		return msgBytes
	}
	return newBytes
}

// CheckMessage 检查消息是否为合法的飞书消息格式
func (l *LarkBot) CheckMessage(msg string) bool {
	if len(msg) == 0 {
		return false
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(msg), &raw); err != nil {
		return false
	}
	msgType, ok := raw["msg_type"].(string)
	if !ok {
		return false
	}
	switch msgType {
	case "text", "post", "image", "share_chat", "interactive":
		return true
	default:
		return false
	}
}

// marshalMessage 将消息包装成飞书接口要求的格式
func marshalMessage(msg interface{}) ([]byte, error) {
	if text, ok := msg.(Text); ok {
		textMsg := message{MsgType: "text", Content: text}
		return marshal(textMsg)
	}
	if post, ok := msg.(POST); ok {
		postMsg := message{MsgType: "post", Content: post}
		return marshal(postMsg)
	}
	if card, ok := msg.(Card); ok {
		cardMsg := message{MsgType: "share_chat", Content: card}
		return marshal(cardMsg)
	}
	if image, ok := msg.(Image); ok {
		imageMsg := message{MsgType: "image", Content: image}
		return marshal(imageMsg)
	}
	if interactive, ok := msg.(CardInteractive); ok {
		interactiveMsg := message{MsgType: "interactive", Content: interactive}
		return marshal(interactiveMsg)
	}
	// 未知类型尝试直接序列化
	return json.Marshal(msg)
}

// marshal 序列化JSON，禁用HTML转义
func marshal(msg interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "")
	err := jsonEncoder.Encode(msg)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
