package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type LarkBot struct {
	Key        string
	WebHookUrl string
	Client     *http.Client
}

const defaultWebHookUrlTemplate = "https://open.feishu.cn/open-apis/bot/v2/hook/%s"

func New(botKey string) *LarkBot {
	bot := LarkBot{
		Key:        botKey,
		WebHookUrl: fmt.Sprintf(defaultWebHookUrlTemplate, botKey),
		Client:     &http.Client{Timeout: 5 * time.Second},
	}
	return &bot
}

func (l *LarkBot) Send(msg interface{}) (bool, error) {
	msgBytes, err := marshalMessage(msg)
	if err != nil {
		return false, err
	}
	return l.SendRaw(msgBytes)
}

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
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	r := Resp{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return false, err
	}
	if r.Code != 200 {
		return false, fmt.Errorf("send message failed, status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return true, nil
}

func (l *LarkBot) CheckMessage(msg string) bool {
	if len(msg) == 0 {
		return false
	}
	interactive := CardInteractive{}
	err := json.Unmarshal([]byte(msg), &interactive)
	if err != nil {
		return false
	}
	text := Text{}
	err = json.Unmarshal([]byte(msg), &text)
	if err != nil {
		return false
	}
	image := Image{}
	err = json.Unmarshal([]byte(msg), &image)
	if err != nil {
		return false
	}
	card := Card{}
	err = json.Unmarshal([]byte(msg), &card)
	if err != nil {
		return false
	}
	post := POST{}
	err = json.Unmarshal([]byte(msg), &post)
	if err != nil {
		return false
	}
	textMsg := Text{}
	err = json.Unmarshal([]byte(msg), &textMsg)
	if err != nil {
		return false
	}
	return true

}

type Resp struct {
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	Code          int    `json:"code"`
	Data          struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

// 将消息包装成企信接口要求的格式，返回 json bytes
func marshalMessage(msg interface{}) ([]byte, error) {
	if text, ok := msg.(Text); ok {
		textMsg := message{MsgType: "text", Content: text}
		return marshal(textMsg)
	}
	if text, ok := msg.(POST); ok {
		textMsg := message{MsgType: "post", Content: text}
		return marshal(textMsg)
	}
	if text, ok := msg.(Card); ok {
		textMsg := message{MsgType: "share_chat", Content: text}
		return marshal(textMsg)
	}
	if text, ok := msg.(Image); ok {
		textMsg := message{MsgType: "image", Content: text}
		return marshal(textMsg)
	}
	if text, ok := msg.(CardInteractive); ok {
		textMsg := message{MsgType: "interactive", Content: text}
		return marshal(textMsg)
	}
	return nil, nil
}

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
