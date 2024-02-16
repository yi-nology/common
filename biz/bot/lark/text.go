package lark

type message struct {
	MsgType string      `json:"msg_type"`
	Content interface{} `json:"content"`
}

type Text struct {
	Text string `json:"text"`
}
