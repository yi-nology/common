package dingtalk

type ActionCardMessage struct {
	MsgType    MsgType     `json:"msgtype"`
	ActionCard *ActionCard `json:"actionCard"`
}

type ActionCard struct {
	Title          string `json:"title"`
	Text           string `json:"text"`
	HideAvatar     string `json:"hideAvatar"`
	BtnOrientation string `json:"btnOrientation"`
	SingleTitle    string `json:"singleTitle"`
	SingleURL      string `json:"singleURL"`
	Btns           []*Btn `json:"btns"`
}

/*** ActionCard ***/
func NewActionCard() *ActionCardMessage {
	return &ActionCardMessage{ACTIONCARD, &ActionCard{}}
}

func (a *ActionCardMessage) SetContent(title, text string) *ActionCardMessage {
	a.ActionCard.Title = title
	a.ActionCard.Text = text
	return a
}

func (a *ActionCardMessage) AddBtn(singleTitle, singleURL string) *ActionCardMessage {
	a.ActionCard.SingleTitle = singleTitle
	a.ActionCard.SingleURL = singleURL
	return a
}

func (a *ActionCardMessage) AddBtns(btns [][]string) *ActionCardMessage {
	for _, item := range btns {
		a.ActionCard.Btns = append(a.ActionCard.Btns, &Btn{item[0], item[1]})
	}
	return a
}

func (a *ActionCardMessage) HideAvatar() *ActionCardMessage {
	a.ActionCard.HideAvatar = "1"
	return a
}

func (a *ActionCardMessage) BtnOrientation() *ActionCardMessage {
	a.ActionCard.BtnOrientation = "0"
	return a
}
