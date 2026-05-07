package dingtalk

type TextMessage struct {
	MsgType MsgType `json:"msgtype"`
	At      *At     `json:"at"`
	Text    *Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
}

/*** Text ***/
func NewText() *TextMessage {
	return &TextMessage{TEXT, &At{}, &Text{}}
}

func (t *TextMessage) SetContent(content string) *TextMessage {
	t.Text.Content = content
	return t
}

func (t *TextMessage) AtMobiles(mobiles []string) *TextMessage {
	t.At.AtMobiles = mobiles
	return t
}

func (t *TextMessage) AtAll() *TextMessage {
	t.At.IsAtAll = true
	return t
}
