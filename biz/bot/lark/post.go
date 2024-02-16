package lark

type POST struct {
	Post struct {
		ZhCn struct {
			Title   string `json:"title"`
			Content [][]struct {
				Tag    string `json:"tag"`
				Text   string `json:"text,omitempty"`
				Href   string `json:"href,omitempty"`
				UserId string `json:"user_id,omitempty"`
			} `json:"content"`
		} `json:"zh_cn"`
	} `json:"post"`
}
