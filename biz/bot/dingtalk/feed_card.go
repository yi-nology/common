package dingtalk

type FeedCardMessage struct {
	MsgType  MsgType   `json:"msgtype"`
	FeedCard *FeedCard `json:"feedCard"`
}

type FeedCard struct {
	Links []*FeedCardLink `json:"links"`
}

type FeedCardLink struct {
	Title      string `json:"title"`
	MessageURL string `json:"messageURL"`
	PicURL     string `json:"picURL"`
}

/*** FeedCard ***/
func NewFeedCard() *FeedCardMessage {
	return &FeedCardMessage{FEEDCARD, &FeedCard{}}
}

func (f *FeedCardMessage) AddCard(title, messageURL, picURL string) *FeedCardMessage {
	f.FeedCard.Links = append(f.FeedCard.Links, &FeedCardLink{title, messageURL, picURL})
	return f
}

func (f *FeedCardMessage) AddCards(cards [][]string) *FeedCardMessage {
	for _, item := range cards {
		f.FeedCard.Links = append(f.FeedCard.Links, &FeedCardLink{item[0], item[1], item[2]})
	}
	return f
}
