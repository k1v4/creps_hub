package service

import (
	"article_service/internal/domain/entity"
	"context"
	"time"
)

type ArticleStorage interface {
	PostArticle(ctx context.Context, name, text string, tags []string, createdAt time.Duration) (int64, error)
	GetArticle(ctx context.Context, id int64) (*entity.Article, error)
}

type articleService struct {
	storage ArticleStorage
}

func NewArticleService(storage ArticleStorage) *articleService {
	return &articleService{storage: storage}
}

func (s *articleService) AddNewArticle(ctx context.Context, name, text string, tags []string, createdAt time.Duration) (int64, error) {
	return s.storage.PostArticle(ctx, name, text, tags, createdAt)
}

func (s *articleService) GetFullArticle(ctx context.Context, id int64) (*entity.Article, error) {
	return s.storage.GetArticle(ctx, id)
}
