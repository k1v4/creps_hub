package entity

type PostArticleRequest struct {
	Title   string `json:"name"`
	Content string `json:"text"`
}

type PostArticleResponse struct {
	Id int `json:"article_id"`
}

type DeleteArticleResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

type PaginatedResponse struct {
	Items  []Article `json:"items"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
