package dingtalk

type MsgType string

const TEXT MsgType = "text"
const LINK MsgType = "link"
const MARKDOWN MsgType = "markdown"
const ACTIONCARD MsgType = "actionCard"
const FEEDCARD MsgType = "feedCard"

type At struct {
	AtMobiles []string `json:"atMobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type Btn struct {
	Title     string `json:"title"`
	ActionURL string `json:"actionURL"`
}
