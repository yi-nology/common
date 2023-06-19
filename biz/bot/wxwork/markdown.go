package wxwork

type markdownMessage struct {
	message
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

func NewMarkdown() *Markdown {
	return &Markdown{}
}

func (m *Markdown) SetMarkdown(s string) *Markdown {
	m.Content = s
	return m
}
