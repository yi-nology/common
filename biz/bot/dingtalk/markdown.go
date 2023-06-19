package dingtalk

type MarkDownMessage struct {
	MsgType  MsgType   `json:"msgtype"`
	At       *At       `json:"at"`
	MarkDown *MarkDown `json:"markdown"`
}

type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

/*** markdown ***/
func NewMarkDown() *MarkDownMessage {
	return &MarkDownMessage{MARKDOWN, &At{}, &MarkDown{}}
}

func (m *MarkDownMessage) SetContent(title, text string) *MarkDownMessage {
	m.MarkDown.Title = title
	m.MarkDown.Text = text
	return m
}

func (m *MarkDownMessage) AtMobiles(mobiles []string) *MarkDownMessage {
	m.At.AtMobiles = mobiles
	return m
}

func (m *MarkDownMessage) AtAll() *MarkDownMessage {
	m.At.IsAtAll = true
	return m
}
