package lanxin
package lanxin

// TextMessage 文本消息
type TextMessage struct {
	MsgType MsgType `json:"msgtype"`
	At      *At     `json:"at,omitempty"`
	Text    *Text   `json:"text"`
}

// Text 文本内容
type Text struct {
	Content string `json:"content"`
}

// NewText 创建文本消息
func NewText() *TextMessage {
	return &TextMessage{MsgType: TEXT, Text: &Text{}}
}

// SetContent 设置文本内容
func (t *TextMessage) SetContent(content string) *TextMessage {
	t.Text.Content = content
	return t
}

// AtMobiles @指定手机号
func (t *TextMessage) AtMobiles(mobiles []string) *TextMessage {
	if t.At == nil {
		t.At = &At{}
	}
	t.At.AtMobiles = mobiles
	return t
}

// AtAll @所有人
func (t *TextMessage) AtAll() *TextMessage {
	if t.At == nil {
		t.At = &At{}
	}
	t.At.IsAtAll = true
	return t
}
