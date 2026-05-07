package wxwork

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

func init() {
}

const (
	defaultWebHookUrlTemplate = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
)

var (
	ErrUnsupportedMessage = errors.New("尚不支持的消息类型")
)

// WxWorkBot 企业微信机器人
type WxWorkBot struct {
	Key        string
	WebHookUrl string
	Client     *http.Client
}

type message struct {
	MsgType string `json:"msgtype"`
}

// New 创建企业微信机器人实例
func New(botKey string) *WxWorkBot {
	bot := WxWorkBot{
		Key:        botKey,
		WebHookUrl: fmt.Sprintf(defaultWebHookUrlTemplate, botKey),
		Client:     &http.Client{Timeout: 5 * time.Second},
	}
	return &bot
}

// Send 发送消息
func (bot *WxWorkBot) Send(msg interface{}) (bool, error) {
	msgBytes, err := marshalMessage(msg)
	if err != nil {
		return false, err
	}
	return bot.SendRaw(msgBytes)
}

// SendRaw 发送原始JSON消息
func (bot *WxWorkBot) SendRaw(msgBytes []byte) (bool, error) {
	webHookUrl := bot.WebHookUrl
	if len(webHookUrl) == 0 {
		webHookUrl = fmt.Sprintf(defaultWebHookUrlTemplate, bot.Key)
	}
	req, err := http.NewRequest(http.MethodPost, webHookUrl, bytes.NewBuffer(msgBytes))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := bot.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var wxWorkResp wxWorkResponse
	err = json.Unmarshal(body, &wxWorkResp)
	if err != nil {
		return false, err
	}
	if wxWorkResp.ErrorCode != 0 && wxWorkResp.ErrorMessage != "" {
		return false, errors.New(string(body))
	}
	return true, nil
}

// CheckMessage 检查消息是否为合法的企业微信消息格式
func (bot *WxWorkBot) CheckMessage(msg string) bool {
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
	case "text", "markdown", "image", "news", "template_card":
		return true
	default:
		return false
	}
}

// marshal 序列化JSON，禁用HTML转义
func marshal(msg interface{}) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buf)
	jsonEncoder.SetEscapeHTML(false)
	jsonEncoder.SetIndent("", "")
	err := jsonEncoder.Encode(msg)
	if err != nil {
		return nil, nil
	}
	return buf.Bytes(), nil
}

// marshalMessage 将消息包装成企业微信接口要求的格式
func marshalMessage(msg interface{}) ([]byte, error) {
	if text, ok := msg.(Text); ok {
		textMsg := textMessage{message: message{MsgType: "text"}, Text: text}
		return marshal(textMsg)
	}
	if textMsg, ok := msg.(textMessage); ok {
		textMsg.MsgType = "text"
		return marshal(textMsg)
	}
	if markdown, ok := msg.(Markdown); ok {
		markdownMsg := markdownMessage{message: message{MsgType: "markdown"}, Markdown: markdown}
		return marshal(markdownMsg)
	}
	if markdownMsg, ok := msg.(markdownMessage); ok {
		markdownMsg.MsgType = "markdown"
		return marshal(markdownMsg)
	}
	if image, ok := msg.(Image); ok {
		imageMsg := imageMessage{message: message{MsgType: "image"}, Image: image}
		return marshal(imageMsg)
	}
	if imageMsg, ok := msg.(imageMessage); ok {
		imageMsg.MsgType = "image"
		return marshal(imageMsg)
	}
	if news, ok := msg.(News); ok {
		newsMsg := newsMessage{message: message{MsgType: "news"}, News: news}
		return marshal(newsMsg)
	}
	if newsMsg, ok := msg.(newsMessage); ok {
		newsMsg.MsgType = "news"
		return marshal(newsMsg)
	}
	if templateCard, ok := msg.(TemplateCard); ok {
		templateCardMsg := templateCardMessage{message: message{MsgType: "template_card"}, TemplateCard: templateCard}
		return marshal(templateCardMsg)
	}
	if templateCardMsg, ok := msg.(templateCardMessage); ok {
		templateCardMsg.MsgType = "template_card"
		return marshal(templateCardMsg)
	}
	return nil, ErrUnsupportedMessage
}
