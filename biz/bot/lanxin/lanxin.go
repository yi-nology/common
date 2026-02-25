package lanxin
package lanxin

// MsgType 消息类型
type MsgType string

const (
	TEXT     MsgType = "text"
	MARKDOWN MsgType = "markdown"
)

// At @功能
type At struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll   bool     `json:"isAtAll,omitempty"`
}
