package bot

import (
	"github.com/yi-nology/common/biz/bot/dingtalk"
	"github.com/yi-nology/common/biz/bot/lanxin"
	"github.com/yi-nology/common/biz/bot/lark"
	"github.com/yi-nology/common/biz/bot/wxwork"
)

// BotOne 机器人统一接口
type BotOne interface {
	Send(interface{}) (bool, error)        // 发送消息
	SendRaw(msgBytes []byte) (bool, error) // 发送原始消息
	CheckMessage(msg string) bool          // 检查消息是否合法
}

// BotType 机器人类型
type BotType string

const (
	WXWork   BotType = "wx"
	Dingtalk BotType = "dd"
	Lark     BotType = "lark"
	Lanxin   BotType = "lanxin"
)

// SecurityType 安全模式
type SecurityType string

const (
	SecurityNone    SecurityType = "none"
	SecuritySign    SecurityType = "sign"
	SecurityKeyword SecurityType = "keyword"
)

// SwitchOne 根据类型创建机器人实例（向后兼容：secret非空时默认签名模式）
func SwitchOne(botType BotType, botKey string, secret string) BotOne {
	switch botType {
	case WXWork:
		return wxwork.New(botKey)
	case Dingtalk:
		return dingtalk.New(botKey).AddSign(secret)
	case Lark:
		return lark.New(botKey).AddSign(secret)
	case Lanxin:
		return lanxin.New(botKey).AddSign(secret)
	default:
		return wxwork.New(botKey)
	}
}

// NewBot 创建机器人实例，支持完整安全模式配置
func NewBot(botType BotType, botKey string, securityType SecurityType, secret string, keyword string) BotOne {
	switch botType {
	case WXWork:
		return wxwork.New(botKey)
	case Dingtalk:
		bot := dingtalk.New(botKey)
		return applyDingtalkSecurity(bot, securityType, secret, keyword)
	case Lark:
		bot := lark.New(botKey)
		return applyLarkSecurity(bot, securityType, secret, keyword)
	case Lanxin:
		bot := lanxin.New(botKey)
		return applyLanxinSecurity(bot, securityType, secret, keyword)
	default:
		return wxwork.New(botKey)
	}
}

func applyDingtalkSecurity(bot *dingtalk.Robot, securityType SecurityType, secret, keyword string) BotOne {
	switch securityType {
	case SecuritySign:
		return bot.AddSign(secret)
	case SecurityKeyword:
		return bot.AddKeyword(keyword)
	default:
		return bot
	}
}

func applyLarkSecurity(bot *lark.LarkBot, securityType SecurityType, secret, keyword string) BotOne {
	switch securityType {
	case SecuritySign:
		return bot.AddSign(secret)
	case SecurityKeyword:
		return bot.AddKeyword(keyword)
	default:
		return bot
	}
}

func applyLanxinSecurity(bot *lanxin.Robot, securityType SecurityType, secret, keyword string) BotOne {
	switch securityType {
	case SecuritySign:
		return bot.AddSign(secret)
	case SecurityKeyword:
		return bot.AddKeyword(keyword)
	default:
		return bot
	}
}
