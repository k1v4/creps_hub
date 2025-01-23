package v1

import (
	"article_service/internal/domain/entity"
	"context"
	"time"
)

const (
	articleURL  = "/article/:article_id"
	articlesURL = "/article"
)

type articleHandler struct {
	articleUseCase ArticleUseCase
}

type ArticleUseCase interface {
	AddNewArticle(ctx context.Context, name, text string, tags []string, createdAt time.Duration) (int64, error)
	GetFullArticle(ctx context.Context, id int64) (*entity.Article, error)
}

func NewArticleHandler(a ArticleUseCase) *articleHandler {
	return &articleHandler{a}
}

func (a *articleHandler) Register() {}
