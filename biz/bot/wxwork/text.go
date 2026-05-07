package wxwork

type textMessage struct {
	message
	Text Text `json:"text"`
}

type Text struct {
	Content             string   `json:"content"`
	MentionedList       []string `json:"mentioned_list"`
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

func NewText() *Text {
	return &Text{}
}

func (t *Text) AddMobile(phone string) *Text {
	t.MentionedMobileList = append(t.MentionedMobileList, phone)
	return t
}

func (t *Text) SetMobileList(phones []string) *Text {
	t.MentionedMobileList = phones
	return t
}

func (t *Text) AddMentioned(mentioned string) *Text {
	t.MentionedList = append(t.MentionedList, mentioned)
	return t
}

func (t *Text) SetMentionedList(mentioned []string) *Text {
	t.MentionedList = mentioned
	return t
}
