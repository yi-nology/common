package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/yi-nology/common/utils/xjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Robot struct {
	Key        string
	RequestUrl string
	Client     *http.Client
}

type Resp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

const (
	WEBHOOK_URL = "https://oapi.dingtalk.com/robot/send?access_token=%s"
)

// New 获取一个 Robot instance
func New(botKey string) *Robot {
	return &Robot{
		Key:        botKey,
		RequestUrl: fmt.Sprintf(WEBHOOK_URL, botKey),
		Client:     &http.Client{Timeout: 5 * time.Second},
	}
}

// AddSign 更新SDK 支持Sign模式
func (r *Robot) AddSign(secret string) *Robot {
	if len(secret) == 0 {
		return r
	}
	microTimestamp := time.Now().UnixNano() / 1e6
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, fmt.Sprintf("%d\n%s", microTimestamp, secret))
	if r.RequestUrl != "" {
		r.RequestUrl = r.RequestUrl + "&timestamp=" + strconv.Itoa(int(microTimestamp)) + "&sign=" + url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	}
	return r
}

// Send 发送notification
func (r *Robot) Send(message interface{}) (bool, error) {
	b, err := json.Marshal(message)
	if err != nil {
		return false, err
	}
	// log.Println(string(b))
	return r.SendRaw(b)
}

func (r *Robot) SendRaw(msgBytes []byte) (bool, error) {
	resp, err := r.Client.Post(r.RequestUrl, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	ret := &Resp{}

	if err := json.Unmarshal(body, &ret); err != nil {
		return false, err
	}

	if ret.ErrCode == 0 {
		return true, nil
	} else {
		return false, errors.New(ret.ErrMsg)
	}
}

func (r *Robot) CheckMessage(msg string) bool {

	if len(msg) == 0 {
		return false
	}
	text := TextMessage{}
	markdown := MarkDownMessage{}
	link := LinkMessage{}
	feedcard := FeedCardMessage{}
	actionCard := ActionCardMessage{}
	if xjson.UnmarshalFromString(msg, &text) == nil {
		return true
	}
	if xjson.UnmarshalFromString(msg, &markdown) == nil {
		return true
	}
	if xjson.UnmarshalFromString(msg, &link) == nil {
		return true
	}
	if xjson.UnmarshalFromString(msg, &feedcard) == nil {
		return true
	}
	if xjson.UnmarshalFromString(msg, &actionCard) == nil {
		return true
	}
	return false
}
