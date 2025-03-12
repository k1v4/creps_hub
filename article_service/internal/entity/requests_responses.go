package entity

type PostArticleRequest struct {
	Title   string `json:"name"`
	Content string `json:"text"`
}

type PostArticleResponse struct {
	Id int `json:"article_id"`
}
