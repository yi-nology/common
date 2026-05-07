package lark

type CardInteractive struct {
	Elements []struct {
		Tag  string `json:"tag"`
		Text struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		} `json:"text,omitempty"`
		Actions []struct {
			Tag  string `json:"tag"`
			Text struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"text"`
			Url   string `json:"url"`
			Type  string `json:"type"`
			Value struct {
			} `json:"value"`
		} `json:"actions,omitempty"`
	} `json:"elements"`
	Header struct {
		Title struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		} `json:"title"`
	} `json:"header"`
}
