package bot

import (
	"github.com/yi-nology/common/biz/bot/dingtalk"
	"github.com/yi-nology/common/biz/bot/wxwork"
)

type BotOne interface {
	Send(interface{}) (bool, error)        // 发送
	SendRaw(msgBytes []byte) (bool, error) // 发送原始消息
	CheckMessage(msg string) bool          // 检查消息是否合法
}

type BotType string

const (
	WXWork   BotType = "wx"
	Dingtalk BotType = "dd"
)

func SwitchOne(botType BotType, botKey string, secret string) BotOne {
	switch botType {
	case WXWork:
		return wxwork.New(botKey)
	case Dingtalk:
		return dingtalk.New(botKey).AddSign(secret)
	default:
		return wxwork.New(botKey)
	}
}
