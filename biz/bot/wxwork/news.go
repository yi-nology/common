package wxwork

type newsMessage struct {
	message
	News News `json:"news"`
}

type News struct {
	Articles []NewsArticle `json:"articles"`
}

func NewNews() *News {
	return &News{}
}

func (n *News) SetArticles(na []NewsArticle) *News {
	n.Articles = na
	return n
}

func (n *News) AddArticle(na NewsArticle) *News {
	n.Articles = append(n.Articles, na)
	return n
}

type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PicURL      string `json:"picurl"`
}
