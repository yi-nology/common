package lanxin

// MarkDownMessage Markdown消息
type MarkDownMessage struct {
	MsgType  MsgType   `json:"msgtype"`
	At       *At       `json:"at,omitempty"`
	MarkDown *MarkDown `json:"markdown"`
}

// MarkDown Markdown内容
type MarkDown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// NewMarkDown 创建Markdown消息
func NewMarkDown() *MarkDownMessage {
	return &MarkDownMessage{MsgType: MARKDOWN, MarkDown: &MarkDown{}}
}

// SetContent 设置Markdown内容
func (m *MarkDownMessage) SetContent(title, text string) *MarkDownMessage {
	m.MarkDown.Title = title
	m.MarkDown.Text = text
	return m
}

// AtMobiles @指定手机号
func (m *MarkDownMessage) AtMobiles(mobiles []string) *MarkDownMessage {
	if m.At == nil {
		m.At = &At{}
	}
	m.At.AtMobiles = mobiles
	return m
}

// AtAll @所有人
func (m *MarkDownMessage) AtAll() *MarkDownMessage {
	if m.At == nil {
		m.At = &At{}
	}
	m.At.IsAtAll = true
	return m
}
