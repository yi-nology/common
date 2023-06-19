package dingtalk

type LinkMessage struct {
	MsgType MsgType `json:"msgtype"`
	Link    *Link   `json:"link"`
}

type Link struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageUrl string `json:"messageUrl"`
	PicUrl     string `json:"picUrl"`
}

/*** Link ***/
func NewLink() *LinkMessage {
	return &LinkMessage{LINK, &Link{}}
}

func (l *LinkMessage) SetContent(title, text, messageUrl, picUrl string) *LinkMessage {
	l.Link.Title = title
	l.Link.Text = text
	l.Link.MessageUrl = messageUrl
	l.Link.PicUrl = picUrl
	return l
}
